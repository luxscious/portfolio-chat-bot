/// <reference types="vite/client" />

interface ViteTypeOptions {
  // By adding this line, you can make the type of ImportMetaEnv strict
  // to disallow unknown keys.
  // strictImportMetaEnv: unknown
}

interface ImportMetaEnv {
  readonly VITE_APP_TITLE: string;
  readonly VITE_API_URL: string;
  readonly VITE_DEVCYCLE_CLIENT_SDK_KEY: string;
  readonly VITE_DEVCYCLE_USER_ID: string;
}

interface ImportMeta {
  readonly env: ImportMetaEnv;
}
