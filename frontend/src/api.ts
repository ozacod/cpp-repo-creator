import type { Library, Category, ProjectConfig } from './types';

const API_BASE = 'http://localhost:8000/api';

export async function fetchLibraries(): Promise<Library[]> {
  const response = await fetch(`${API_BASE}/libraries`);
  if (!response.ok) throw new Error('Failed to fetch libraries');
  const data = await response.json();
  return data.libraries;
}

export async function fetchCategories(): Promise<Category[]> {
  const response = await fetch(`${API_BASE}/categories`);
  if (!response.ok) throw new Error('Failed to fetch categories');
  const data = await response.json();
  return data.categories;
}

export async function searchLibraries(query: string): Promise<Library[]> {
  const response = await fetch(`${API_BASE}/search?q=${encodeURIComponent(query)}`);
  if (!response.ok) throw new Error('Search failed');
  const data = await response.json();
  return data.results;
}

export async function previewCMake(config: ProjectConfig): Promise<string> {
  const response = await fetch(`${API_BASE}/preview`, {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    body: JSON.stringify(config),
  });
  
  if (!response.ok) throw new Error('Preview failed');
  const data = await response.json();
  return data.cmake_content;
}

export async function generateProject(config: ProjectConfig): Promise<Blob> {
  const response = await fetch(`${API_BASE}/generate`, {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    body: JSON.stringify(config),
  });
  
  if (!response.ok) {
    const error = await response.json();
    throw new Error(error.detail || 'Generation failed');
  }
  
  return response.blob();
}

export async function reloadRecipes(): Promise<void> {
  const response = await fetch(`${API_BASE}/reload-recipes`, {
    method: 'POST',
  });
  if (!response.ok) throw new Error('Failed to reload recipes');
}
