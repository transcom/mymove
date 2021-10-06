UPDATE duty_stations
    SET provides_services_counseling = TRUE
    WHERE name in (
        'NAVSTA Newport',
        'NSB New London',
        'Portsmouth Naval Shipyard',
        'JB Anacostia–Bolling',
        'Station New York',
        'Camp Lejeune',
        'JBSA Fort Sam Houston',
        'Sector Mobile',
        'JB Elmendorf-Richardson',
        'NAVSUP FLC Puget Sound',
        'Eielson AFB',
        'Ellington Field ANGB',
        'Griffiss AFB'
);