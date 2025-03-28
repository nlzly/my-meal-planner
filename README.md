# My Meal Planner

A modern web application for planning your weekly meals, built with React, TypeScript, and Vite.

## Features

- 🗓️ Weekly meal planning grid
- 🍳 Add, edit, and delete meals
- 👨‍🍳 Assign chefs to meals
- 📝 Add descriptions to meals
- 🎯 Drag and drop meals between slots
- ⌨️ Keyboard shortcuts for copy/paste
- 💾 Local storage persistence
- 🔒 Authentication support (mock implementation)

## Tech Stack

- React 18
- TypeScript
- Vite
- CSS Modules
- Local Storage API

## Getting Started

### Prerequisites

- Node.js (v14 or higher)
- npm or yarn

### Installation

1. Clone the repository:
```bash
git clone https://github.com/yourusername/my-meal-planner.git
cd my-meal-planner
```

2. Install dependencies:
```bash
cd client-vite
npm install
```

3. Start the development server:
```bash
npm run dev
```

The application will be available at `http://localhost:3000`.

### Building for Production

To create a production build:

```bash
npm run build
```

The built files will be in the `dist` directory.

To preview the production build:

```bash
npm run preview
```

## Project Structure

```
my-meal-planner/
├── client-vite/
│   ├── src/
│   │   ├── components/     # React components
│   │   ├── services/       # API and local storage services
│   │   ├── types/         # TypeScript type definitions
│   │   ├── App.tsx        # Main application component
│   │   ├── main.tsx       # Application entry point
│   │   ├── App.css        # Global styles
│   │   └── index.css      # Base styles
│   ├── public/            # Static assets
│   ├── index.html         # HTML template
│   ├── package.json       # Dependencies and scripts
│   ├── tsconfig.json      # TypeScript configuration
│   ├── vite.config.ts     # Vite configuration
│   └── README.md          # Project documentation
└── README.md              # Root README
```

## Development

### Available Scripts

- `npm run dev` - Start the development server
- `npm run build` - Create a production build
- `npm run preview` - Preview the production build
- `npm run lint` - Run ESLint

### Code Style

The project uses ESLint for code linting and TypeScript for type checking. The configuration can be found in:
- `.eslintrc.cjs`
- `tsconfig.json`

## Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

This project is licensed under the MIT License - see the LICENSE file for details.
