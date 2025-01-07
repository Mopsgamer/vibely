import type { Config } from "tailwindcss";
const config: Config = {
    content: [
        "./web/templates/**/*.html",
    ],
    theme: {
        extend: {
            backgroundImage: {
                "mesh-gradient": "url(/static/assets/image-mesh-gradient.png)",
            },
            transitionDuration: {
                "1100": "1100ms",
                "1200": "1200ms",
                "1300": "1300ms",
                "1400": "1400ms",
            },
            spacing: {
                "page": "1.25rem",
            },
            colors: {
                "neutral-0": "var(--sl-color-neutral-0)",
                "neutral-50": "var(--sl-color-neutral-50)",
                "neutral-100": "var(--sl-color-neutral-100)",
                "neutral-200": "var(--sl-color-neutral-200)",
                "neutral-300": "var(--sl-color-neutral-300)",
                "neutral-400": "var(--sl-color-neutral-400)",
                "neutral-500": "var(--sl-color-neutral-500)",
                "neutral-600": "var(--sl-color-neutral-600)",
                "neutral-700": "var(--sl-color-neutral-700)",
                "neutral-800": "var(--sl-color-neutral-800)",
                "neutral-900": "var(--sl-color-neutral-900)",
                "neutral-950": "var(--sl-color-neutral-950)",
                "neutral-1000": "var(--sl-color-neutral-1000)",

                "primary-50": "var(--sl-color-primary-50)",
                "primary-100": "var(--sl-color-primary-100)",
                "primary-200": "var(--sl-color-primary-200)",
                "primary-300": "var(--sl-color-primary-300)",
                "primary-400": "var(--sl-color-primary-400)",
                "primary-500": "var(--sl-color-primary-500)",
                "primary-600": "var(--sl-color-primary-600)",
                "primary-700": "var(--sl-color-primary-700)",
                "primary-800": "var(--sl-color-primary-800)",
                "primary-900": "var(--sl-color-primary-900)",
                "primary-950": "var(--sl-color-primary-950)",

                "danger-50": "var(--sl-color-danger-50)",
                "danger-100": "var(--sl-color-danger-100)",
                "danger-200": "var(--sl-color-danger-200)",
                "danger-300": "var(--sl-color-danger-300)",
                "danger-400": "var(--sl-color-danger-400)",
                "danger-500": "var(--sl-color-danger-500)",
                "danger-600": "var(--sl-color-danger-600)",
                "danger-700": "var(--sl-color-danger-700)",
                "danger-800": "var(--sl-color-danger-800)",
                "danger-900": "var(--sl-color-danger-900)",
                "danger-950": "var(--sl-color-danger-950)",

                "warning-50": "var(--sl-color-warning-50)",
                "warning-100": "var(--sl-color-warning-100)",
                "warning-200": "var(--sl-color-warning-200)",
                "warning-300": "var(--sl-color-warning-300)",
                "warning-400": "var(--sl-color-warning-400)",
                "warning-500": "var(--sl-color-warning-500)",
                "warning-600": "var(--sl-color-warning-600)",
                "warning-700": "var(--sl-color-warning-700)",
                "warning-800": "var(--sl-color-warning-800)",
                "warning-900": "var(--sl-color-warning-900)",
                "warning-950": "var(--sl-color-warning-950)",

                "success-50": "var(--sl-color-success-50)",
                "success-100": "var(--sl-color-success-100)",
                "success-200": "var(--sl-color-success-200)",
                "success-300": "var(--sl-color-success-300)",
                "success-400": "var(--sl-color-success-400)",
                "success-500": "var(--sl-color-success-500)",
                "success-600": "var(--sl-color-success-600)",
                "success-700": "var(--sl-color-success-700)",
                "success-800": "var(--sl-color-success-800)",
                "success-900": "var(--sl-color-success-900)",
                "success-950": "var(--sl-color-success-950)",
            },
        },
    },
};
export default config;
