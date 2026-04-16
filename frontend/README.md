# Calendar Booking Frontend

A React-based frontend for the Calendar Booking System, built with TypeScript, Vite, and Mantine UI.

## Tech Stack

- **TypeScript** - Type-safe development
- **Vite** - Fast build tool and dev server
- **React SPA** - Single Page Application
- **Mantine UI** - Comprehensive UI component library
- **Axios** - HTTP client for API integration
- **React Router** - Client-side routing
- **Day.js** - Date/time handling and formatting
- **Prism** - API mocking for development

## Features

### Guest (Public) Features
- View available event types with descriptions and durations
- Select an event type and browse available time slots on a calendar
- Book a time slot with name and email validation
- Email format validation
- Real-time feedback and error handling

### Owner Features
- Create, edit, and delete event types
- Define event type properties (name, description, duration)
- View all upcoming bookings in a paginated list
- Filter and sort bookings
- Cancel bookings
- Pagination support for large datasets

## Getting Started

### Prerequisites

- Node.js 18+
- npm or yarn

### Installation

```bash
npm install
```

### Development

Start the development server:

```bash
npm run dev
```

The app will be available at `http://localhost:5173`

### API Mocking

For development without a backend, use Prism mock server:

```bash
# From the project root (one level up)
./frontend/mock-api.sh
```

Or manually:

```bash
npx @stoplight/prism-cli mock -h 0.0.0.0 -p 4010 ../typespec/tsp-output/schema/openapi.yaml
```

The mock API will be available at `http://localhost:4010`

### Environment Variables

- `.env.development` - Uses mock API at `http://localhost:4010`
- `.env.production` - Uses real backend at `http://localhost:8080`

To override:

```bash
# Create a local env file
cp .env.development .env.local
# Edit .env.local with your preferred API URL
```

### Build

Build for production:

```bash
npm run build
```

The built files will be in `dist/`

### Preview Production Build

```bash
npm run preview
```

## Project Structure

```
frontend/
├── src/
│   ├── api/
│   │   └── client.ts          # API client with Axios
│   ├── components/
│   │   ├── layout/
│   │   │   └── Header.tsx     # App header component
│   │   └── owner/
│   │       ├── EventTypeManagement.tsx   # CRUD for event types
│   │       └── BookingsList.tsx          # View/cancel bookings
│   ├── pages/
│   │   ├── guest/
│   │   │   ├── GuestHome.tsx   # Public event types listing
│   │   │   └── BookingPage.tsx # Calendar and booking form
│   │   ├── owner/
│   │   │   └── OwnerDashboard.tsx  # Owner main page
│   │   └── NotFound.tsx        # 404 page
│   ├── types/
│   │   └── api.ts              # TypeScript type definitions
│   ├── App.tsx                 # Main app component with routing
│   ├── main.tsx                # App entry point
│   └── index.css               # Global styles
├── .env.development            # Dev environment variables
├── .env.production             # Production environment variables
├── mock-api.sh                 # Script to start Prism mock server
└── package.json
```

## API Integration

The frontend integrates with the backend API following the TypeSpec contract:

### Owner API (`/api`)
- `POST /api/event-types` - Create event type
- `GET /api/event-types` - List event types (paginated)
- `GET /api/event-types/:id` - Get event type details
- `PUT /api/event-types/:id` - Update event type
- `DELETE /api/event-types/:id` - Delete event type
- `GET /api/bookings` - List all bookings (paginated)
- `GET /api/bookings/:id` - Get booking details
- `DELETE /api/bookings/:id` - Cancel booking
- `GET /api/slots` - List all time slots (paginated)

### Guest API (`/api/public`)
- `GET /api/public/event-types` - List available event types (paginated)
- `GET /api/public/event-types/:id` - Get event type details
- `GET /api/public/slots` - List available public slots (paginated)
- `POST /api/public/bookings` - Create a booking

## Validation

### Input Validation
- **Email**: Valid email format required
- **Duration**: Minimum 5 minutes, maximum 1440 minutes (24 hours)
- **Name**: Required, 1-100 characters
- **Description**: Required, 1-500 characters

### Business Rules
- No overlapping bookings allowed (enforced by backend)
- Cannot book past time slots
- All times use UTC with timezone support

## Error Handling

The app provides comprehensive error handling:

- **400 Bad Request**: Validation errors with field-specific messages
- **404 Not Found**: Resource not found
- **409 Conflict**: Booking time conflicts
- **Network errors**: User-friendly error messages
- Toast notifications for all operations

## Routing

- `/` - Guest home page (public event types)
- `/book/:eventTypeId` - Booking page for specific event type
- `/owner` - Owner dashboard (event types and bookings management)
- `*` - 404 page for unknown routes

## Browser Support

The app supports modern browsers:
- Chrome 90+
- Firefox 88+
- Safari 14+
- Edge 90+

## Scripts

```bash
npm run dev          # Start dev server
npm run build        # Build for production
npm run preview      # Preview production build
npm run lint         # Run ESLint
```

## Acceptance Criteria

✅ Fully functional UI adhering to the API spec
✅ Covers both user roles (Owner and Guest)
✅ Enforces booking rules (no overlaps, validation)
✅ Handles validation errors correctly
✅ Handles HTTP errors correctly
✅ Supports pagination for all list views
✅ Input validation (email format, duration > 0)
✅ Calendar-based slot selection
✅ Responsive design with Mantine UI
