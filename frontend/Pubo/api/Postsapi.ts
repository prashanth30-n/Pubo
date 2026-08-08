import { PostDraft,PostStatus } from "@/types/Post";
import { Platform } from 'react-native';

// Set EXPO_PUBLIC_API_URL in your .env, e.g. EXPO_PUBLIC_API_URL=https://api.yourapp.com/api
const API_BASE = process.env.EXPO_PUBLIC_API_URL ?? 'http://localhost:8080/api';

export interface PostResponse {
  id: string;
  content: string;
  platforms: string[];
  images: string[];
  status: PostStatus;
  scheduledAt: string | null;
  createdAt: string;
  updatedAt: string;
}
export interface SavedPost { id: string; content: string; status: 'draft' | 'posted' | 'failed'; mediaIds: string[]; platforms: string[]; createdAt: string; media: { id: string; storageURL: string; contentType: string }[]; }
export interface PostPage { data: SavedPost[]; limit: number; offset: number; total: number; }

async function request<T>(path: string, options: RequestInit = {}): Promise<T> {
  const res = await fetch(`${API_BASE}${path}`, {
    headers: { 'Content-Type': 'application/json' },
    ...options,
  });
  if (!res.ok) {
    const body = await res.text().catch(() => '');
    throw new Error(`Request failed (${res.status}): ${body}`);
  }
  if (res.status === 204) return undefined as T;
  return res.json();
}

type PostPayload = Pick<PostDraft, 'content' | 'platforms'> & { images: string[] };

export const postsApi = {
  saveDraft(token: string, payload: { content: string; mediaIds: string[] }) {
    return request<SavedPost>('/v1/posts/drafts', { method: 'POST', headers: { 'Content-Type': 'application/json', Authorization: `Bearer ${token}` }, body: JSON.stringify(payload) });
  },
  list(token: string, status?: 'draft' | 'posted' | 'failed', limit = 10, offset = 0) {
    const query = `?status=${status ?? ''}&limit=${limit}&offset=${offset}`;
    return request<PostPage>(`/v1/posts${query}`, { headers: { Authorization: `Bearer ${token}` } });
  },
  publish(token: string, payload: { content: string; platforms: string[]; mediaIds: string[] }) {
    return request<{ results: { platform: string; status: string; externalPostId?: string; error?: string }[] }>('/v1/posts/publish', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json', Authorization: `Bearer ${token}` },
      body: JSON.stringify(payload),
    });
  },

  async uploadImage(token: string, uri: string, name: string, type: string, webFile?: File) {
    const form = new FormData();
    // Browsers require a Blob/File. Native Expo expects the { uri, name, type }
    // object, so retain that path for iOS/Android.
    if (Platform.OS === 'web') {
      const file: Blob = (webFile ?? await fetch(uri).then((response) => response.blob())) as Blob;
      form.append('file', file, name);
    } else {
      form.append('file', { uri, name, type } as any);
    }
    form.append('folder', 'posts');
    const response = await fetch(`${API_BASE}/v1/media/upload`, { method: 'POST', headers: { Authorization: `Bearer ${token}` }, body: form });
    if (!response.ok) throw new Error(await response.text());
    return response.json() as Promise<{ id: string; storageURL: string }>;
  },
  createDraft(payload: PostPayload) {
    return request<PostResponse>('/posts', {
      method: 'POST',
      body: JSON.stringify({ ...payload, status: 'draft' }),
    });
  },

  updatePost(id: string, payload: Partial<PostPayload>) {
    return request<PostResponse>(`/posts/${id}`, {
      method: 'PUT',
      body: JSON.stringify(payload),
    });
  },

  schedulePost(id: string, scheduledAt: string) {
    return request<PostResponse>(`/posts/${id}/schedule`, {
      method: 'POST',
      body: JSON.stringify({ scheduledAt }),
    });
  },

  publishNow(id: string) {
    return request<PostResponse>(`/posts/${id}/publish`, { method: 'POST' });
  },

  deletePost(id: string) {
    return request<void>(`/posts/${id}`, { method: 'DELETE' });
  },

  listPosts(status?: PostStatus) {
    const query = status ? `?status=${status}` : '';
    return request<PostResponse[]>(`/posts${query}`);
  },
};
