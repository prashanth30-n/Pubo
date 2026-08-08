import { useCallback, useState } from 'react';
import { postsApi } from '@/api/Postsapi';
import { PlatformId,PostImage,PostStatus } from '@/types/Post';

interface Options {
  initialContent?: string;
  initialPlatforms?: PlatformId[];
}

export function usePostComposer(options: Options = {}) {
  const [postId, setPostId] = useState<string | undefined>();
  const [content, setContent] = useState(options.initialContent ?? '');
  const [platforms, setPlatforms] = useState<PlatformId[]>(options.initialPlatforms ?? []);
  const [images, setImages] = useState<PostImage[]>([]);
  const [status, setStatus] = useState<PostStatus>('draft');
  const [scheduledAt, setScheduledAt] = useState<string | null>(null);
  const [isSaving, setIsSaving] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const togglePlatform = useCallback((id: PlatformId) => {
    setPlatforms((prev) => (prev.includes(id) ? prev.filter((p) => p !== id) : [...prev, id]));
  }, []);

  const addImage = useCallback((image: PostImage) => {
    setImages((prev) => [...prev, image]);
  }, []);

  const removeImage = useCallback((id: string) => {
    setImages((prev) => prev.filter((img) => img.id !== id));
  }, []);

  // Creates the post on first save, updates it on every save after that —
  // so draft/schedule/publish all share one "make sure it's saved" step.
  const persist = useCallback(async () => {
    const payload = { content, platforms, images: images.map((img) => img.url) };
    if (postId) {
      return postsApi.updatePost(postId, payload);
    }
    const created = await postsApi.createDraft(payload);
    setPostId(created.id);
    return created;
  }, [content, platforms, images, postId]);

  const saveDraft = useCallback(async (token: string) => {
    setIsSaving(true);
    setError(null);
    try {
      const saved = await postsApi.saveDraft(token, { content, mediaIds: images.map((image) => image.id) });
      setPostId(saved.id);
      setStatus('draft');
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Could not save draft');
    } finally {
      setIsSaving(false);
    }
  }, [content, images]);

  const schedule = useCallback(
    async (isoTimestamp: string) => {
      setIsSaving(true);
      setError(null);
      try {
        const saved = await persist();
        await postsApi.schedulePost(saved.id, isoTimestamp);
        setScheduledAt(isoTimestamp);
        setStatus('scheduled');
      } catch (err) {
        setError(err instanceof Error ? err.message : 'Could not schedule post');
      } finally {
        setIsSaving(false);
      }
    },
    [persist]
  );

  const publishNow = useCallback(async (token: string) => {
    setIsSaving(true);
    setError(null);
    try {
      const result = await postsApi.publish(token, { content, platforms, mediaIds: images.map((image) => image.id) });
      if (result.results.some((platform) => platform.status === 'failed')) throw new Error(result.results.filter((platform) => platform.status === 'failed').map((platform) => `${platform.platform}: ${platform.error}`).join('\n'));
      setStatus('published');
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Could not publish post');
    } finally {
      setIsSaving(false);
    }
  }, [content, platforms, images]);

  const canSubmit = content.trim().length > 0 && platforms.length > 0;

  return {
    content,
    setContent,
    platforms,
    togglePlatform,
    images,
    addImage,
    removeImage,
    status,
    scheduledAt,
    isSaving,
    error,
    canSubmit,
    saveDraft,
    schedule,
    publishNow,
  };
}
