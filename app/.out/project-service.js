import window from "./window.js";
export class ProjectService {
    url = null;
    constructor() {
        this.setUrl();
    }
    setUrl() {
        this.url = window.vars.SERVER_URL;
    }
    /**
     * Get All Projects
     */
    async getAllProjects() {
        const res = await fetch(`${this.url}/api/projects`);
        if (!res.ok)
            throw new Error('Failed to fetch projects');
        return res.json();
    }
    /**
     * Get Project
     */
    async getProject(id) {
        const res = await fetch(`${this.url}/api/projects/${id}`);
        if (!res.ok)
            throw new Error('Failed to fetch projects');
        return res.json();
    }
    /**
     * Create Project
     */
    async createProject(data) {
        const res = await fetch(`${this.url}/api/projects`, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify(data)
        });
        if (!res.ok) {
            throw new Error('Failed to create project');
        }
        return res.json();
    }
    /**
     * Update Project
     */
    async updateProject(id, data) {
        const res = await fetch(`${this.url}/api/projects/${id}`, {
            method: 'PUT',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify(data)
        });
        if (!res.ok) {
            throw new Error('Failed to update project');
        }
        return res.json();
    }
    /**
     * Delete Project
     */
    async deleteProject(id) {
        const res = await fetch(`${this.url}/api/projects/${id}`, {
            method: 'DELETE'
        });
        if (!res.ok) {
            throw new Error('Failed to delete project');
        }
        return res.json();
    }
}
