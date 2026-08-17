import { useCallback, useState } from 'react';
import { listComments } from '../api/comment';
import type { Comment } from '../types';

export function useCommentStore() {
  const [comments, setComments] = useState<Comment[]>([]);
  const [loading, setLoading] = useState(false);

  const refresh = useCallback(async (articleId: number) => {
    setLoading(true);
    try {
      setComments(await listComments(articleId));
    } finally {
      setLoading(false);
    }
  }, []);

  return { comments, loading, refresh };
}
