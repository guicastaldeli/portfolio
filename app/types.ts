export interface Project {
    id: number;
    name: string;
    desc: string;
    repo: string;
    createdAt: string;
    updatedAt: string;
    media: Media[]
    links: Link[];
}

export interface Media {
    id: number;
    projectId: number;
    type: string;
    url: string;
}

export interface Link {
    id: number;
    projectId: number;
    name: string;
    url: string;
}

export interface CreateProjectRequest {
    name: string;
    desc: string;
    repo: string;
    photos: string[];
    videos: string[];
    links: { name: string; url: string }[];
}

export interface WebSocketMessage {
    type: string;
    channel: string;
    data: any;
}