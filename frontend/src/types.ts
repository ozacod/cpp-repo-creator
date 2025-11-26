export type Library = {
  id: string;
  name: string;
  description: string;
  category: string;
  tags: string[];
  github_url: string;
  fetch_content: string;
  link_libraries: string[];
  header_only: boolean;
  cpp_standard: number;
  alternatives: string[];
};

export type Category = {
  id: string;
  name: string;
  icon: string;
  description: string;
};

export type ProjectConfig = {
  project_name: string;
  cpp_standard: number;
  library_ids: string[];
  include_tests: boolean;
};
