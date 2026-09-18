/// <reference types="vite/client" />

// Sem isso, o Vite tipa qualquer chave de import.meta.env como `any`
// (ImportMetaEnvFallbackKey em vite/types/importMeta.d.ts) — declarando as
// variáveis explicitamente, `import.meta.env.VITE_API_BASE_URL` vira `string`
// de verdade em vez de escapar do modo estrito.
interface ImportMetaEnv {
  readonly VITE_API_BASE_URL: string;
}

interface ImportMeta {
  readonly env: ImportMetaEnv;
}
