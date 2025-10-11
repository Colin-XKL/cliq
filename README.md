# cliQ

**cliQ** /klɪk/, similar to "click"
cliQ is a lightweight tool designed to transform complex CLI commands into intuitive and user-friendly graphical user interfaces (GUIs).

Users only need to define a command template with variable placeholders, and the corresponding GUI interface will be automatically generated. Parameters can be filled in through clicking and selecting, eliminating the need to memorize command syntax and significantly improving user experience and operational efficiency.

## Key Features
- Command template definition: Users can input CLI commands to generate corresponding dynamic form interfaces.
- Template import/export: Supports exporting templates (`.cliqfile.yaml`) as files or importing from others, facilitating team collaboration.
- Multiple input component support: Supports file selectors, number input boxes, dropdown selectors, checkboxes, etc., adapting to different parameter types.
- Cross-platform support: Supports Windows, macOS, and Linux platforms
- Template marketplace: Users can upload/download templates for common tools (such as ImageMagick, ffmpeg, pngquant, etc.) to build a shared ecosystem.

## Tech Stack
- Built with the wails framework for cross-platform support
- Frontend uses Vue3 with TypeScript
- Uses PrimeVue V4 component library, Vite to build the frontend project, with TailwindCSS for some styling