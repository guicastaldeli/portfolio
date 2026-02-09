import { ProjectService } from "./project-service.js";
import { GetProjectHandler } from "./get-project-handler.js";
import type { Project } from "./types.js";

export class Main {
    private projectService: ProjectService;
    private projectHandler: GetProjectHandler;
    private currentProjects: Project[] = [];

    constructor() {
        this.projectService = new ProjectService();
        this.projectHandler = new GetProjectHandler();
        this.init();
    }

    private async init(): Promise<void> {
        await this.projectHandler.connect();
        this.setupHandlers();
        await this.loadProjects();
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

    private renderProjects(): void {
        console.log('Rendering projects...');
        const container = document.getElementById('projects-container');
        console.log('Container found:', container);
        
        if (!container) {
            console.error('projects-container not found!');
            return;
        }
        
        container.innerHTML = '';

        if (!this.currentProjects || this.currentProjects.length === 0) {
            console.log('No projects to display');
            container.innerHTML = '<p>No projects yet.</p>';
            return;
        }

        console.log(`Rendering ${this.currentProjects.length} projects`);
        
        this.currentProjects.forEach(project => {
            const photos = project.media.filter(m => m.type === 'photo');
            const videos = project.media.filter(m => m.type === 'video');

            const projectContainer = document.createElement('div');
            projectContainer.className = 'project-container';
            let mediaHtml = '';

            if (photos.length > 0) {
                mediaHtml += '<div class="project-photos">';
                photos.forEach(photo => {
                    mediaHtml += `
                        <div class="photo-item">
                            <img src="${this.escapeHtml(photo.url)}" 
                                alt="Project photo" 
                                loading="lazy"
                                onerror="this.style.display='none'">
                        </div>
                    `;
                });
                mediaHtml += '</div>';
            }

            if (videos.length > 0) {
                mediaHtml += '<div class="project-videos">';
                videos.forEach(video => {
                    if (this.isVideoUrl(video.url)) {
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

            let linksHtml = '';
            if (project.links.length > 0) {
                linksHtml += '<div class="project-links">';
                project.links.forEach(link => {
                    linksHtml += `
                        <a href="${this.escapeHtml(link.url)}" 
                        target="_blank" 
                        class="project-link">
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
                                    <a href="${this.escapeHtml(project.repo)}" target="_blank">
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

    private isVideoUrl(url: string): boolean {
        const videoExtensions = ['.mp4', '.webm', '.ogg', '.mov', '.avi', '.mkv'];
        const lowerUrl = url.toLowerCase();
        return videoExtensions.some(ext => lowerUrl.endsWith(ext));
    }

    private truncate(
        text: string,
        charsPerLine: number,
        maxLines: number
    ): string {
        const maxChars = charsPerLine * maxLines;
        const sliced = text.slice(0, maxChars);
        const lines: string[] = [];

        for (let i = 0; i < sliced.length; i += charsPerLine) {
            lines.push(sliced.slice(i, i + charsPerLine));
        }

        const needsEllipsis = text.length > maxChars;
        if (needsEllipsis) lines[lines.length - 1] += '...';

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