import type { Project, CreateProjectRequest } from "./types.js";
import window from "./window.js";

export class ProjectService {
    private url: string | null = null;

    constructor() {
        this.setUrl();
    }

    private setUrl(): void {
        this.url = window.vars.SERVER_URL;
    }

    /**
     * Get All Projects
     */
    public async getAllProjects(): Promise<Project[]> {
        const res = await fetch(`${this.url}/api/projects`);
        if(!res.ok) throw new Error('Failed to fetch projects');
        return res.json();
    }

    /**
     * Get Project
     */
    public async getProject(id: number): Promise<Project> {
        const res = await fetch(`${this.url}/api/projects/${id}`);
        if(!res.ok) throw new Error('Failed to fetch projects');
        return res.json();
    }

    /**
     * Create Project
     */
    public async createProject(data: CreateProjectRequest): Promise<{ id: number; message: string }> {
        const res = await fetch(`${this.url}/api/projects`, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify(data)
        });
        if(!res.ok) {
            throw new Error('Failed to create project');
        }
        return res.json();
    }

    /**
     * Update Project
     */
    public async updateProject(id: number, data: CreateProjectRequest): Promise<{ message: string }> {
        const res = await fetch(`${this.url}/api/projects/${id}`, {
            method: 'PUT',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify(data)
        });
        if(!res.ok) {
            throw new Error('Failed to update project');
        }
        return res.json();
    }

    /**
     * Delete Project
     */
    public async deleteProject(id: number): Promise<{ message: string }> {
        const res = await fetch(`${this.url}/api/projects/${id}`, {
            method: 'DELETE'
        });
        if(!res.ok) {
            throw new Error('Failed to delete project');
        }
        return res.json();
    }
}