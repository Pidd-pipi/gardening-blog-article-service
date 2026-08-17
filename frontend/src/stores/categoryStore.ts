import { useCallback, useState } from 'react';
import { getCategoryTree } from '../api/category';
import type { Category } from '../types';

let cache: Category[] = [];

export function useCategoryStore() {
  const [tree, setTree] = useState<Category[]>(cache);

  const refresh = useCallback(async (): Promise<Category[]> => {
    const data = await getCategoryTree();
    cache = data;
    setTree(data);
    return data;
  }, []);

  return { tree, refresh };
}
