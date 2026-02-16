import { ProjectService } from "./project-service.js";
import { GetProjectHandler } from "./get-project-handler.js";
import type { Project } from "./types.js";
import window from "./window.js";

export class Main {
    private projectService: ProjectService;
    private projectHandler: GetProjectHandler;
    private currentProjects: Project[] = [];

    constructor() {
        this.projectService = new ProjectService();
        this.projectHandler = new GetProjectHandler();

        
        window.init();
        this.init();

        window.vars.APP_ENV = 'prod';
    }

    private async init(): Promise<void> {
        this.connect();
        await this.projectHandler.connect();
        this.setupHandlers();
        await this.loadProjects();
        this.createModal();
    }

    private connect(): WebSocket {
        const ws = new WebSocket(window.vars.SERVER_WS);
        return ws;
    }

    /**
     * Setup WebSocket Handlers
     */
    private setupHandlers(): void {
        /* On Project Created */
        this.projectHandler.setOnProjectCreated((data) => {
            console.log('Project created', data);
            this.loadProjects();
        });
        /* On Project Updated */
        this.projectHandler.setOnProjectUpdated((data) => {
            console.log('Project updated', data);
            this.loadProjects();
        });
        /* On Project Deleted */
        this.projectHandler.setOnProjectDeleted((data) => {
            console.log('Project deleted', data);
            this.loadProjects();
        });
    }

    private async loadProjects(): Promise<void> {
        try {
            console.log('Loading projects...');
            this.currentProjects = await this.projectService.getAllProjects();
            console.log('Projects loaded:', this.currentProjects);
            this.renderProjects();
        } catch (err) {
            console.error('Failed to load projects', err);
        }
    }

    /**
     * Create Modal
     */
    private createModal(): void {
        const modal = document.createElement('div');
        modal.id = 'project-modal';
        modal.className = 'project-modal hidden';
        modal.innerHTML = `
            <div class="modal-overlay"></div>
            <div class="modal-content">
                <button class="modal-close">&times;</button>
                <div id="modal-project-details"></div>
            </div>
        `;
        document.body.appendChild(modal);

        modal.querySelector('.modal-overlay')?.addEventListener('click', () => {
            this.closeModal();
        });
        modal.querySelector('.modal-close')?.addEventListener('click', () => {
            this.closeModal();
        });

        document.addEventListener('keydown', (e) => {
            if(e.key === 'Escape') {
                this.closeModal();
            }
        });
    }

    /**
     * Show Project Details
     */
    private showProjectDetails(project: Project): void {
        const modal = document.getElementById('project-modal');
        const detailsContainer = document.getElementById('modal-project-details');
        
        if(!modal || !detailsContainer) return;

        const photos = project.media.filter(m => m.type === 'photo');
        const videos = project.media.filter(m => m.type === 'video');

        let mediaHtml = '';

        if(photos.length > 0) {
            mediaHtml += '<div class="modal-photos">';
            photos.forEach(photo => {
                mediaHtml += `
                    <div class="modal-photo-item">
                        <img src="${this.escapeHtml(photo.url)}" 
                            alt="Project photo" 
                            loading="lazy"
                            onerror="this.style.display='none'">
                    </div>
                `;
            });
            mediaHtml += '</div>';
        }
        if(videos.length > 0) {
            mediaHtml += '<div class="modal-videos">';
            videos.forEach(video => {
                const id = this.getVideoId(video.url);
                if(id) {
                    const thumbnailUrl = `https://img.youtube.com/vi/${id}/hqdefault.jpg`;
                    mediaHtml += `
                        <div class="modal-video-item">
                            <a href="${this.escapeHtml(video.url)}" target="_blank" class="video-thumbnail-link">
                                <img src="${thumbnailUrl}" 
                                    alt="Video thumbnail" 
                                    class="video-thumbnail"
                                    onerror="this.src='https://img.youtube.com/vi/${id}/hqdefault.jpg'"
                                >
                                <div class="play-button-overlay">▶</div>
                            </a>
                        </div>
                    `;
                } else if(this.isVideoUrl(video.url)) {
                    mediaHtml += `
                        <div class="modal-video-item">
                            <video controls width="100%">
                                <source src="${this.escapeHtml(video.url)}" type="video/mp4">
                                Your browser does not support the video tag.
                            </video>
                        </div>
                    `;
                }
            });
            mediaHtml += '</div>';
        }
        let linksHtml = '';
        if(project.links.length > 0) {
            linksHtml += '<div class="modal-links"><h4>Links</h4><ul>';
            project.links.forEach(link => {
                linksHtml += `
                    <li>
                        <a href="${this.escapeHtml(link.url)}" 
                           target="_blank" 
                           class="modal-link">
                            ${this.escapeHtml(link.name)}
                        </a>
                    </li>
                `;
            });
            linksHtml += '</ul></div>';
        }

        detailsContainer.innerHTML = `
            <div class="modal-header">
                <h2>${this.truncate(this.escapeHtml(project.name), 80)}</h2>
            </div>
            <div class="modal-body">
                ${mediaHtml}
                ${linksHtml}

                ${project.repo ? `
                    <div class="modal-repo">
                        <h4>Repository</h4>
                        <a href="${this.escapeHtml(project.repo)}" target="_blank">
                            ${this.escapeHtml(project.repo)}
                        </a>
                    </div>
                ` : ''}
                
                <div class="modal-description">
                    <h4>Description</h4>
                    <p>${this.truncate(this.escapeHtml(project.desc), 80)}</p>
                </div>
            </div>
        `;

        modal.classList.remove('hidden');
        document.body.style.overflow = 'hidden';
    }

    /**
     * Close Modal
     */
    private closeModal(): void {
        const modal = document.getElementById('project-modal');
        if(modal) {
            modal.classList.add('hidden');
            document.body.style.overflow = '';
        }
    }

    private renderProjects(): void {
        console.log('Rendering projects...');
        const container = document.getElementById('projects-container');
        
        if(!container) {
            console.error('projects-container not found!');
            return;
        }
        
        container.innerHTML = '';

        if(!this.currentProjects || this.currentProjects.length === 0) {
            console.log('No projects to display');
            container.innerHTML = '<p>No projects yet.</p>';
            return;
        }

        console.log(`Rendering ${this.currentProjects.length} projects`);
        
        this.currentProjects.forEach(project => {
            const photos = project.media.filter(m => m.type === 'photo');
            const videos = project.media.filter(m => m.type === 'video');
            
            const MAX_PREVIEW_MEDIA = 3;
            const allMedia = [...photos, ...videos];
            const previewMedia = allMedia.slice(0, MAX_PREVIEW_MEDIA);
            const hasMoreMedia = allMedia.length > MAX_PREVIEW_MEDIA;

            const projectContainer = document.createElement('div');
            projectContainer.className = 'project-container';
            
            projectContainer.addEventListener('click', () => {
                this.showProjectDetails(project);
            });

            let mediaHtml = '';
            
            const previewPhotos = previewMedia.filter(m => m.type === 'photo');
            const previewVideos = previewMedia.filter(m => m.type === 'video');

            if(previewPhotos.length > 0) {
                mediaHtml += '<div class="project-photos">';
                previewPhotos.forEach(photo => {
                    mediaHtml += `
                        <div class="photo-item">
                            <img src="${this.escapeHtml(photo.url)}" 
                                alt="Project photo" 
                                loading="lazy"
                                onerror="this.style.display='none'"
                            >
                        </div>
                    `;
                });
                mediaHtml += '</div>';
            }
            if(previewVideos.length > 0) {
                mediaHtml += '<div class="project-videos">';
                previewVideos.forEach(video => {
                    const id = this.getVideoId(video.url);
                    if(id) {
                        const thumbnailUrl = `https://img.youtube.com/vi/${id}/hqdefault.jpg`;
                        mediaHtml += `
                            <div class="video-item">
                                <div class="video-thumbnail-link">
                                    <img src="${thumbnailUrl}" 
                                        alt="Video thumbnail" 
                                        class="video-thumbnail"
                                        onerror="this.src='https://img.youtube.com/vi/${id}/hqdefault.jpg'"
                                    >
                                    <div class="play-button-overlay">▶</div>
                                </div>
                            </div>
                        `;
                    } else if(this.isVideoUrl(video.url)) {
                        mediaHtml += `
                            <div class="video-item">
                                <video controls width="200">
                                    <source src="${this.escapeHtml(video.url)}" type="video/mp4">
                                    Your browser does not support the video tag.
                                </video>
                            </div>
                        `;
                    }
                });
                mediaHtml += '</div>';
            }
            
            if(hasMoreMedia) {
                mediaHtml += `
                    <div class="more-media-indicator">
                        +${allMedia.length - MAX_PREVIEW_MEDIA} more
                    </div>
                `;
            }
            
            let linksHtml = '';
            if(project.links.length > 0) {
                linksHtml += '<div class="project-links">';
                project.links.forEach(link => {
                    linksHtml += `
                        <a href="${this.escapeHtml(link.url)}" 
                        target="_blank" 
                        class="project-link"
                        onclick="event.stopPropagation()">
                        ${this.escapeHtml(link.name)}
                        </a>
                    `;
                });
                linksHtml += '</div>';
            }

            projectContainer.innerHTML = `
                <div class="project-main">
                    <div id="project-main-content">
                        ${mediaHtml}
                        <div id="project-info">
                            <h3 id="project-name">${this.escapeHtml(project.name)}</h3>
                            ${project.repo ? `
                                <div class="project-repo">
                                    <strong>Repository:</strong> 
                                    <a href="${this.escapeHtml(project.repo)}" 
                                       target="_blank"
                                       onclick="event.stopPropagation()">
                                        ${this.escapeHtml(project.repo)}
                                    </a>
                                </div>
                            ` : ''}
                            <p class="project-description">${this.truncate(this.escapeHtml(project.desc), 25, 3)}</p>
                            ${linksHtml}
                        </div>
                    </div>
                </div>
            `;

            container.appendChild(projectContainer);
        });
    }

    private getVideoId(url: string): string | null {
        const patterns = [
            /(?:youtube\.com\/watch\?v=|youtu\.be\/|youtube\.com\/embed\/)([^&\n?#]+)/,
            /youtube\.com\/shorts\/([^&\n?#]+)/
        ];
        
        for (const pattern of patterns) {
            const match = url.match(pattern);
            if (match && match[1]) {
                return match[1];
            }
        }
        return null;
    }

    private isVideoUrl(url: string): boolean {
        const videoExtensions = ['.mp4', '.webm', '.ogg', '.mov', '.avi', '.mkv'];
        const lowerUrl = url.toLowerCase();
        return videoExtensions.some(ext => lowerUrl.endsWith(ext));
    }

    private truncate(
        text: string,
        charsPerLine: number,
        maxLines?: number
    ): string {
        const lines: string[] = [];

        if(maxLines) {
            const maxChars = charsPerLine * maxLines;
            const sliced = text.slice(0, maxChars);

            for(let i = 0; i < sliced.length; i += charsPerLine) {
                lines.push(sliced.slice(i, i + charsPerLine));
            }

            const needsEllipsis = text.length > maxChars;
            if(needsEllipsis) lines[lines.length - 1] += '...';
        } else {
            for(let i = 0; i < text.length; i += charsPerLine) {
                lines.push(text.slice(i, i + charsPerLine));
            }
        }

        return lines.join('<br>');
    }

    private escapeHtml(text: string): string {
        const div = document.createElement('div');
        div.textContent = text;
        return div.innerHTML;
    }

    /**
     * Cleanup
     */
    public cleanup(): void {
        this.projectHandler.disconnect();
    }
}

document.addEventListener('DOMContentLoaded', () => {
    console.log('DOM ready, creating Main');
    new Main();
});