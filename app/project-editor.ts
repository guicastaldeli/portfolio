import type { Project, CreateProjectRequest } from "./types.js";
import { ProjectService } from "./project-service.js";
import { GetProjectHandler } from "./get-project-handler.js";
import { Main } from "./server/main.js";

export class ProjectEditor {
    private main: Main;
    private projectService: ProjectService;
    private projectHandler: GetProjectHandler;
    
    private currentProjects: Project[] = [];
    private editingProjectId: number | null = null;

    private el: HTMLSpanElement | null = null;

    constructor(main: Main) {
        this.main = main;

        this.projectService = new ProjectService();
        this.projectHandler = new GetProjectHandler();
        
        this.init();
        this.setLink();
    }

    private setLink() {
        this.el = document.querySelector('.main #project-editor #pe-link');
        if(this.el) {
            this.el.textContent = 'Open Editor';
            (this.el as HTMLAnchorElement).href = '/editor';
        }
    }

    /**
     * 
     * Init
     * 
     */
    private async init(): Promise<void> {
        await this.projectHandler.connect();
        this.setupHandlers();

        await this.loadProjects();

        this.setupEventListeners();
    }

    /**
     * Setup Handlers
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

    /**
     * 
     * Load Projects
     * 
     */
    private async loadProjects(): Promise<void> {
        try {
            this.currentProjects = await this.projectService.getAllProjects();
            this.renderProjects();
        } catch(err) {
            console.error('Failed to load projects', err);
        }
    }

    /**
     * 
     * Render Projects
     * 
     */
    private renderProjects(): void {
        const container = document.getElementById('projects-container');
        if(!container) return;
        container.innerHTML = '';

        if(!this.currentProjects || this.currentProjects.length === 0) {
            container.innerHTML = '<p>No projects yet. Create one!</p>';
            return;
        }

        this.currentProjects.forEach(project => {
            const photos = project.media.filter(m => m.type === 'photo');
            const videos = project.media.filter(m => m.type === 'video');
            
            const projectContainer = document.createElement('div');
            projectContainer.className = 'project-container';
            let mediaHtml = '';
            
            if(photos.length > 0) {
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
            if(videos.length > 0) {
                mediaHtml += '<div class="project-videos">';
                videos.forEach(video => {
                    if(this.isVideoUrl(video.url)) {
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
            if(project.links.length > 0) {
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
                <div class="project-metadata">
                    <div id="project-metadata-content">
                        <small>Created: ${new Date(project.createdAt).toLocaleDateString()}</small>
                        <small>Updated: ${new Date(project.updatedAt).toLocaleDateString()}</small>
                    </div>
                </div>
                ${mediaHtml}
                <div class="project-info">
                    <div id="project-info-content">
                        <h3>${this.escapeHtml(project.name)}</h3>
                        ${project.repo ? `
                            <div class="project-repo">
                                <strong>Repository:</strong> 
                                <a href="${this.escapeHtml(project.repo)}" target="_blank">
                                    ${this.escapeHtml(project.repo)}
                                </a>
                            </div>
                        ` : ''}
                        <p class="project-description">${this.escapeHtml(project.desc)}</p>
                        ${linksHtml}
                    </div>
                </div>
                <div class="project-actions">
                    <button class="edit-btn" data-id="${project.id}">Edit</button>
                    <button class="delete-btn" data-id="${project.id}">Delete</button>
                </div>
            `;

            // Edit Button
            projectContainer.querySelector('.edit-btn')?.addEventListener('click', () => {
                this.editProject(project.id);
            });

            // Delete Button
            projectContainer.querySelector('.delete-btn')?.addEventListener('click', () => {
                this.deleteProject(project.id);
            });

            container.appendChild(projectContainer);
        });
    }

    private isVideoUrl(url: string): boolean {
        const videoExtensions = ['.mp4', '.webm', '.ogg', '.mov', '.avi', '.mkv'];
        const lowerUrl = url.toLowerCase();
        return videoExtensions.some(ext => lowerUrl.endsWith(ext));
    }

    /**
     * 
     * Edit Project
     * 
     */
    private async editProject(id: number): Promise<void> {
        try {
            const project = await this.projectService.getProject(id);
            this.showForm(project);
        } catch(err) {
            console.error('Failed to load project:', err);
        }
    }

    /**
     * 
     * Delete Project
     * 
     */
    private async deleteProject(id: number): Promise<void> {
        try {
            await this.projectService.deleteProject(id);
            console.log('Project deleted!');
            await this.loadProjects();
        } catch(err) {
            console.error('Failed to delete project!', err);
        }
    }

    /**
     * 
     * 
     * --- Setup Event Listeners ---
     * 
     * 
     */
    private setupEventListeners(): void {
        // Create New Project Button
        document.getElementById('create-new-btn')?.addEventListener('click', () => {
            this.showForm();
        });

        // Form Submit
        document.getElementById('edit-form')?.addEventListener('submit', (e) => {
            e.preventDefault();
            this.handleSubmit();
        });

        // Cancel Button
        document.getElementById('cancel-btn')?.addEventListener('click', () => {
            this.hideForm();
        });

        // Add Photo Button
        document.getElementById('add-photo-btn')?.addEventListener('click', () => {
            this.addPhotoInput();
        });

        // Add Video Button
        document.getElementById('add-video-btn')?.addEventListener('click', () => {
            this.addVideoInput();
        });

        // Add Link Button
        document.getElementById('add-link-btn')?.addEventListener('click', () => {
            this.addLinkInput();
        });
    }

    /**
     * 
     * 
     * --- Form ---
     * 
     * 
     */
    private showForm(project?: Project): void {
        const form = document.getElementById('project-form');
        const formTitle = document.getElementById('form-title');
        const projectList = document.getElementById('project-list');
        if(!form || !formTitle || !projectList) return;

        if(project) {
            formTitle.textContent = 'Edit Project';
            this.editingProjectId = project.id;
            this.populateForm(project);
        } else {
            formTitle.textContent = 'Create Project';
            this.editingProjectId = null;
            this.resetForm();
        }

        form.classList.remove('hidden');
        projectList.classList.add('hidden');
    }

    private hideForm(): void {
        const form = document.getElementById('project-form');
        const projectList = document.getElementById('project-list');

        if(!form || !projectList) return;

        form.classList.add('hidden');
        projectList.classList.remove('hidden');
        this.resetForm();
    }

    private populateForm(project: Project): void {
        (document.getElementById('project-name') as HTMLInputElement).value = project.name;
        (document.getElementById('project-desc') as HTMLTextAreaElement).value = project.desc;
        (document.getElementById('project-repo') as HTMLInputElement).value = project.repo || '';

        const photosContainer = document.getElementById('photos-container');
        if(photosContainer) {
            photosContainer.innerHTML = '';
            const photos = project.media.filter(m => m.type === 'photo');
            if(photos.length === 0) {
                this.addPhotoInput();
            } else {
                photos.forEach(photo => {
                    this.addPhotoInput(photo.url);
                });
            }
        }

        const videosContainer = document.getElementById('videos-container');
        if(videosContainer) {
            videosContainer.innerHTML = '';
            const videos = project.media.filter(m => m.type === 'video');
            if(videos.length === 0) {
                this.addVideoInput();
            } else {
                videos.forEach(video => {
                    this.addVideoInput(video.url);
                });
            }
        }

        const linksContainer = document.getElementById('links-container');
        if(linksContainer) {
            linksContainer.innerHTML = '';
            if(project.links.length === 0) {
                this.addLinkInput();
            } else {
                project.links.forEach(link => {
                    this.addLinkInput(link.name, link.url);
                });
            }
        }
    }

    private resetForm(): void {
        (document.getElementById('edit-form') as HTMLFormElement).reset();
        
        const photosContainer = document.getElementById('photos-container');
        if(photosContainer) {
            photosContainer.innerHTML = '';
            this.addPhotoInput();
        }

        const videosContainer = document.getElementById('videos-container');
        if(videosContainer) {
            videosContainer.innerHTML = '';
            this.addVideoInput();
        }

        const linksContainer = document.getElementById('links-container');
        if(linksContainer) {
            linksContainer.innerHTML = '';
            this.addLinkInput();
        }
    }

    private addPhotoInput(value: string = ''): void {
        const container = document.getElementById('photos-container');
        if(!container) return;

        const input = document.createElement('input');
        input.type = 'url';
        input.className = 'photo-input';
        input.placeholder = 'Photo URL';
        input.value = value;
        container.appendChild(input);
    }

    private addVideoInput(value: string = ''): void {
        const container = document.getElementById('videos-container');
        if(!container) return;

        const input = document.createElement('input');
        input.type = 'url';
        input.className = 'video-input';
        input.placeholder = 'Video URL';
        input.value = value;
        container.appendChild(input);
    }

    private addLinkInput(name: string = '', url: string = ''): void {
        const container = document.getElementById('links-container');
        if(!container) return;

        const group = document.createElement('div');
        group.className = 'link-input-group';
        group.innerHTML = `
            <input type="text" class="link-name" placeholder="Link name" value="${this.escapeHtml(name)}">
            <input type="url" class="link-url" placeholder="URL" value="${this.escapeHtml(url)}">
            <button type="button" class="remove-link-btn">×</button>
        `;

        group.querySelector('.remove-link-btn')?.addEventListener('click', () => {
            group.remove();
        });

        container.appendChild(group);
    }

    private async handleSubmit(): Promise<void> {
        const name = (document.getElementById('project-name') as HTMLInputElement).value;
        const desc = (document.getElementById('project-desc') as HTMLTextAreaElement).value;
        const repo = (document.getElementById('project-repo') as HTMLInputElement).value;

        const photos = Array.from(document.querySelectorAll('.photo-input'))
            .map(input => (input as HTMLInputElement).value)
            .filter(val => val.trim() !== '');

        const videos = Array.from(document.querySelectorAll('.video-input'))
            .map(input => (input as HTMLInputElement).value)
            .filter(val => val.trim() !== '');

        const linkGroups = Array.from(document.querySelectorAll('.link-input-group'));
        const links = linkGroups.map(group => ({
            name: (group.querySelector('.link-name') as HTMLInputElement).value,
            url: (group.querySelector('.link-url') as HTMLInputElement).value,
        })).filter(link => link.name.trim() !== '' && link.url.trim() !== '');

        const data: CreateProjectRequest = {
            name,
            desc,
            repo,
            photos,
            videos,
            links,
        };

        try {
            if(this.editingProjectId) {
                await this.projectService.updateProject(this.editingProjectId, data);
                console.log('Project updated successfully!');
            } else {
                await this.projectService.createProject(data);
                console.log('Project created successfully!');
            }
            this.hideForm();
            await this.loadProjects();
        } catch (error) {
            console.error('Failed to save project:', error);
            console.log('Failed to save project');
        }
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