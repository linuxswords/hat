# Webapp Sites and Navigation Proposal

## Site Structure

### Main Navigation
- **Dashboard** - Overview of current tournaments and quick actions
- **Tournaments** - Tournament management
- **Archers** - Archer management
- **Bow Classes** - Bow class management and configuration
- **Handicaps** - Handicap set management
- **Reports** - Generated reports and rankings

## Page Hierarchy

### 1. Dashboard (`/`)
- Welcome page with overview cards
- Recent tournaments
- Quick actions (Create Tournament, Add Scores)
- Upcoming tournament alerts

### 2. Tournaments (`/tournaments`)
#### 2.1 Tournament List (`/tournaments`)
- List all tournaments (past, current, upcoming)
- Filter by date, status
- Actions: Create, Edit, Delete, View Details

#### 2.2 Create Tournament (`/tournaments/new`)
- Form to enter tournament data (name, date, location)
- Select handicap set
- Save and continue to archer assignment

#### 2.3 Tournament Details (`/tournaments/{id}`)
- Tournament information display
- Tabs for:
  - **Details** - Basic tournament info, edit options
  - **Participants** - Assigned archers list, add/remove archers
  - **Scores** - Score entry and management
  - **Rankings** - Live ranking list with adjusted scores
  - **Reports** - Generate PDF rankings, send emails

#### 2.4 Score Entry (`/tournaments/{id}/scores`)
- Table/form for entering raw scores per archer
- Real-time calculation of adjusted scores
- Bulk score entry options
- Save and update functionality

#### 2.5 Rankings (`/tournaments/{id}/rankings`)
- Sortable table ranked by adjusted score
- Columns: Rank, Archer Name, Bowclass, Raw Score, Handicap Factor, Adjusted Score
- Export to PDF button
- Print-friendly view

### 3. Archers (`/archers`)
#### 3.1 Archer List (`/archers`)
- List all registered archers
- Search and filter by name, bowclass
- Actions: Add, Edit, Delete, View Profile

#### 3.2 Archer Profile (`/archers/{id}`)
- Archer details (name, bowclass, contact info)
- Tournament history
- Performance statistics

### 4. Bow Classes (`/bowclasses`)
#### 4.1 Bow Class List (`/bowclasses`)
- List all available bow classes
- Display ID, name, age group, gender, and bow type
- Search and filter by bow type, age group, or gender
- Actions: Add, Edit, Delete, View Details

#### 4.2 Create Bow Class (`/bowclasses/new`)
- Form to create new bow class
- Enter ID (e.g., "AMLB"), name, description
- Select age group, gender, and bow type
- Validation to ensure unique ID

#### 4.3 Bow Class Details (`/bowclasses/{id}`)
- View bow class information and description
- List archers currently using this bow class
- Show handicap factors across different handicap sets
- Edit bow class properties
- Usage statistics in tournaments

### 5. Handicaps (`/handicaps`)
#### 5.1 Handicap Sets (`/handicaps`)
- List available handicap sets
- Create, edit, delete handicap sets
- Set default handicap set

#### 5.2 Handicap Details (`/handicaps/{id}`)
- View handicap factors by bowclass
- Edit handicap values
- Usage history in tournaments

### 6. Reports (`/reports`)
- Tournament results archive
- Historical rankings
- Archer performance reports
- Export options (PDF, CSV)

## Navigation Flow Examples

### Tournament Creation Flow:
1. Dashboard → "Create Tournament" button
2. `/tournaments/new` → Enter basic details
3. `/tournaments/{id}` → Assign archers
4. `/tournaments/{id}/scores` → Enter scores
5. `/tournaments/{id}/rankings` → View results
6. Generate PDF and send emails

### Score Management Flow:
1. Dashboard → Select active tournament
2. `/tournaments/{id}` → "Scores" tab
3. `/tournaments/{id}/scores` → Enter/update scores
4. `/tournaments/{id}/rankings` → View updated rankings

## User Interface Considerations

### Responsive Design
- Mobile-friendly for score entry on tablets
- Desktop optimized for tournament management
- Print-friendly ranking views

### Key Features
- Breadcrumb navigation
- Quick action buttons
- Real-time score calculations
- Confirmation dialogs for destructive actions
- Success/error notifications
- Auto-save for score entries

### Accessibility
- Keyboard navigation support
- Screen reader compatible
- High contrast mode option
- Large text options for score entry
