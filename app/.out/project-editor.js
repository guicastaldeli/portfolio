import { ProjectService } from "./project-service.js";
import { GetProjectHandler } from "./get-project-handler.js";
export class ProjectEditor {
    main;
    projectService;
    projectHandler;
    currentProjects = [];
    editingProjectId = null;
    el = null;
    constructor(main) {
        this.main = main;
        this.projectService = new ProjectService();
        this.projectHandler = new GetProjectHandler();
    }
    setLink() {
        this.el = document.querySelector('.main #project-editor #pe-link');
        if (this.el) {
            this.el.textContent = 'Open Editor';
            this.el.href = '/editor';
        }
    }
    /**
     *
     * Init
     *
     */
    async init() {
        await this.set();
        this.setLink();
    }
    async set() {
        await this.projectHandler.connect();
        this.setupHandlers();
        await this.loadProjects();
        this.setupEventListeners();
    }
    /**
     * Setup Handlers
     */
    setupHandlers() {
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
    async loadProjects() {
        try {
            this.currentProjects = await this.projectService.getAllProjects();
            this.renderProjects();
        }
        catch (err) {
            console.error('Failed to load projects', err);
        }
    }
    /**
     *
     * Render Projects
     *
     */
    renderProjects() {
        const container = document.getElementById('projects-container');
        if (!container)
            return;
        container.innerHTML = '';
        if (!this.currentProjects || this.currentProjects.length === 0) {
            container.innerHTML = '<p>No projects yet. Create one!</p>';
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
            let mediaHtml = '';
            const previewPhotos = previewMedia.filter(m => m.type === 'photo');
            const previewVideos = previewMedia.filter(m => m.type === 'video');
            if (previewPhotos.length > 0) {
                mediaHtml += '<div class="project-photos">';
                previewPhotos.forEach(photo => {
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
            if (previewVideos.length > 0) {
                mediaHtml += '<div class="project-videos">';
                previewVideos.forEach(video => {
                    const id = this.getVideoId(video.url);
                    if (id) {
                        const thumbnailUrl = `https://img.youtube.com/vi/${id}/hqdefault.jpg`;
                        mediaHtml += `
                            <div class="video-item">
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
                    }
                    else if (this.isVideoUrl(video.url)) {
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
            if (hasMoreMedia) {
                mediaHtml += `
                    <div class="more-media-indicator">
                        +${allMedia.length - MAX_PREVIEW_MEDIA} more
                    </div>
                `;
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
                <div class="project-metadata">
                    <div id="project-metadata-content">
                        <small>Created: ${new Date(project.createdAt).toLocaleDateString()}</small>
                        <small>Updated: ${new Date(project.updatedAt).toLocaleDateString()}</small>
                    </div>
                </div>
                <div class="project-main">
                    <div id="project-main-content">
                        ${mediaHtml}
                        <div id="project-info">
                            <h3 class="project-name">${this.escapeHtml(project.name)}</h3>
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
    getVideoId(url) {
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
    isVideoUrl(url) {
        const videoExtensions = ['.mp4', '.webm', '.ogg', '.mov', '.avi', '.mkv'];
        const lowerUrl = url.toLowerCase();
        return videoExtensions.some(ext => lowerUrl.endsWith(ext));
    }
    truncate(text, charsPerLine, maxLines) {
        const maxChars = charsPerLine * maxLines;
        const sliced = text.slice(0, maxChars);
        const lines = [];
        for (let i = 0; i < sliced.length; i += charsPerLine) {
            lines.push(sliced.slice(i, i + charsPerLine));
        }
        const needsEllipsis = text.length > maxChars;
        if (needsEllipsis)
            lines[lines.length - 1] += '...';
        return lines.join('<br>');
    }
    /**
     *
     * Edit Project
     *
     */
    async editProject(id) {
        try {
            const project = await this.projectService.getProject(id);
            this.showForm(project);
        }
        catch (err) {
            console.error('Failed to load project:', err);
        }
    }
    /**
     *
     * Delete Project
     *
     */
    async deleteProject(id) {
        try {
            await this.projectService.deleteProject(id);
            console.log('Project deleted!');
            await this.loadProjects();
        }
        catch (err) {
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
    setupEventListeners() {
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
    showForm(project) {
        const form = document.getElementById('project-form');
        const formTitle = document.getElementById('form-title');
        const projectList = document.getElementById('project-list');
        const header = document.getElementById('header');
        if (!form || !formTitle || !projectList)
            return;
        if (project) {
            formTitle.textContent = 'Edit Project';
            this.editingProjectId = project.id;
            this.populateForm(project);
        }
        else {
            formTitle.textContent = 'Create Project';
            this.editingProjectId = null;
            this.resetForm();
        }
        form.classList.remove('hidden');
        projectList.classList.add('hidden');
        if (header)
            header.style.display = 'none';
    }
    hideForm() {
        const form = document.getElementById('project-form');
        const projectList = document.getElementById('project-list');
        const header = document.getElementById('header');
        if (!form || !projectList)
            return;
        form.classList.add('hidden');
        projectList.classList.remove('hidden');
        this.resetForm();
        if (header)
            header.style.display = 'flex';
    }
    populateForm(project) {
        document.getElementById('project-name').value = project.name;
        document.getElementById('project-desc').value = project.desc;
        document.getElementById('project-repo').value = project.repo || '';
        const photosContainer = document.getElementById('photos-container');
        if (photosContainer) {
            photosContainer.innerHTML = '';
            const photos = project.media.filter(m => m.type === 'photo');
            if (photos.length === 0) {
                this.addPhotoInput();
            }
            else {
                photos.forEach(photo => {
                    this.addPhotoInput(photo.url);
                });
            }
        }
        const videosContainer = document.getElementById('videos-container');
        if (videosContainer) {
            videosContainer.innerHTML = '';
            const videos = project.media.filter(m => m.type === 'video');
            if (videos.length === 0) {
                this.addVideoInput();
            }
            else {
                videos.forEach(video => {
                    this.addVideoInput(video.url);
                });
            }
        }
        const linksContainer = document.getElementById('links-container');
        if (linksContainer) {
            linksContainer.innerHTML = '';
            if (project.links.length === 0) {
                this.addLinkInput();
            }
            else {
                project.links.forEach(link => {
                    this.addLinkInput(link.name, link.url);
                });
            }
        }
    }
    resetForm() {
        document.getElementById('edit-form').reset();
        const photosContainer = document.getElementById('photos-container');
        if (photosContainer) {
            photosContainer.innerHTML = '';
            this.addPhotoInput();
        }
        const videosContainer = document.getElementById('videos-container');
        if (videosContainer) {
            videosContainer.innerHTML = '';
            this.addVideoInput();
        }
        const linksContainer = document.getElementById('links-container');
        if (linksContainer) {
            linksContainer.innerHTML = '';
            this.addLinkInput();
        }
    }
    addPhotoInput(value = '') {
        const container = document.getElementById('photos-container');
        if (!container)
            return;
        const input = document.createElement('input');
        input.type = 'url';
        input.className = 'photo-input';
        input.placeholder = 'Photo URL';
        input.value = value;
        container.appendChild(input);
    }
    addVideoInput(value = '') {
        const container = document.getElementById('videos-container');
        if (!container)
            return;
        const input = document.createElement('input');
        input.type = 'url';
        input.className = 'video-input';
        input.placeholder = 'Video URL';
        input.value = value;
        container.appendChild(input);
    }
    addLinkInput(name = '', url = '') {
        const container = document.getElementById('links-container');
        if (!container)
            return;
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
    async handleSubmit() {
        const name = document.getElementById('project-name').value;
        const desc = document.getElementById('project-desc').value;
        const repo = document.getElementById('project-repo').value;
        const photos = Array.from(document.querySelectorAll('.photo-input'))
            .map(input => input.value)
            .filter(val => val.trim() !== '');
        const videos = Array.from(document.querySelectorAll('.video-input'))
            .map(input => input.value)
            .filter(val => val.trim() !== '');
        const linkGroups = Array.from(document.querySelectorAll('.link-input-group'));
        const links = linkGroups.map(group => ({
            name: group.querySelector('.link-name').value,
            url: group.querySelector('.link-url').value,
        })).filter(link => link.name.trim() !== '' && link.url.trim() !== '');
        const data = {
            name,
            desc,
            repo,
            photos,
            videos,
            links,
        };
        try {
            if (this.editingProjectId) {
                await this.projectService.updateProject(this.editingProjectId, data);
                console.log('Project updated successfully!');
            }
            else {
                await this.projectService.createProject(data);
                console.log('Project created successfully!');
            }
            this.hideForm();
            await this.loadProjects();
        }
        catch (error) {
            console.error('Failed to save project:', error);
            console.log('Failed to save project');
        }
    }
    escapeHtml(text) {
        const div = document.createElement('div');
        div.textContent = text;
        return div.innerHTML;
    }
    /**
     * Cleanup
     */
    cleanup() {
        this.projectHandler.disconnect();
    }
}
