export type LibraryOption = {
  id: string;
  name: string;
  description: string;
  type: 'boolean' | 'string' | 'choice' | 'integer';
  default: any;
  choices?: string[];
  cmake_var?: string;
  cmake_define?: string;
  affects_link?: boolean;
  link_libraries_when_enabled?: string[];
};

export type FetchContent = {
  repository: string;
  tag: string;
  source_subdir?: string;
};

export type Library = {
  id: string;
  name: string;
  description: string;
  category: string;
  github_url: string;
  cpp_standard: number;
  header_only: boolean;
  tags: string[];
  alternatives: string[];
  fetch_content?: FetchContent;
  link_libraries: string[];
  options: LibraryOption[];
  cmake_pre?: string;
  cmake_post?: string;
  system_package?: boolean;
  find_package_name?: string;
};

export type Category = {
  id: string;
  name: string;
  icon: string;
  description: string;
};

export type LibrarySelection = {
  library_id: string;
  options: Record<string, any>;
};

export type ProjectConfig = {
  project_name: string;
  cpp_standard: number;
  libraries: LibrarySelection[];
  include_tests: boolean;
  build_shared: boolean;
};
