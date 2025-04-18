export const ORDERS_TYPE = {
  PERMANENT_CHANGE_OF_STATION: 'PERMANENT_CHANGE_OF_STATION',
  LOCAL_MOVE: 'LOCAL_MOVE',
  RETIREMENT: 'RETIREMENT',
  SEPARATION: 'SEPARATION',
  TEMPORARY_DUTY: 'TEMPORARY_DUTY',
  EARLY_RETURN_OF_DEPENDENTS: 'EARLY_RETURN_OF_DEPENDENTS',
  STUDENT_TRAVEL: 'STUDENT_TRAVEL',
};

export const SPECIAL_ORDERS_TYPES = {
  WOUNDED_WARRIOR: 'Wounded Warrior',
  BLUEBARK: 'BLUEBARK',
  SAFETY: 'Safety',
  SAFETY_NON_LABEL: 'SAFETY',
};

export const CHECK_SPECIAL_ORDERS_TYPES = (ordersType) => {
  return ['BLUEBARK', 'WOUNDED_WARRIOR', 'SAFETY'].includes(ordersType);
};

export const ORDERS_TYPE_OPTIONS = {
  PERMANENT_CHANGE_OF_STATION: 'Permanent Change Of Station (PCS)',
  LOCAL_MOVE: 'Local Move',
  RETIREMENT: 'Retirement',
  SEPARATION: 'Separation',
  WOUNDED_WARRIOR: 'Wounded Warrior',
  BLUEBARK: 'BLUEBARK',
  TEMPORARY_DUTY: 'Temporary Duty (TDY)',
  EARLY_RETURN_OF_DEPENDENTS: 'Early Return of Dependents',
  STUDENT_TRAVEL: 'Student Travel',
};

export const ORDERS_TYPE_DETAILS = {
  HHG_PERMITTED: 'HHG_PERMITTED',
  PCS_TDY: 'PCS_TDY',
  HHG_RESTRICTED_PROHIBITED: 'HHG_RESTRICTED_PROHIBITED',
  HHG_RESTRICTED_AREA: 'HHG_RESTRICTED_AREA',
  INSTRUCTION_20_WEEKS: 'INSTRUCTION_20_WEEKS',
  HHG_PROHIBITED_20_WEEKS: 'HHG_PROHIBITED_20_WEEKS',
  DELAYED_APPROVAL: 'DELAYED_APPROVAL',
};

export const ORDERS_TYPE_DETAILS_OPTIONS = {
  HHG_PERMITTED: 'Shipment of HHG Permitted',
  PCS_TDY: 'PCS with TDY Enroute',
  HHG_RESTRICTED_PROHIBITED: 'Shipment of HHG Restricted or Prohibited',
  HHG_RESTRICTED_AREA: 'HHG Restricted Area-HHG Prohibited',
  INSTRUCTION_20_WEEKS: 'Course of Instruction 20 Weeks or More',
  HHG_PROHIBITED_20_WEEKS: 'Shipment of HHG Prohibited but Authorized within 20 weeks',
  DELAYED_APPROVAL: 'Delayed Approval 20 Weeks or More',
};

export const ORDERS_PAY_GRADE_TYPE = {
  E_1: 'E_1',
  E_2: 'E_2',
  E_3: 'E_3',
  E_4: 'E_4',
  E_5: 'E_5',
  E_6: 'E_6',
  E_7: 'E_7',
  E_8: 'E_8',
  E_9: 'E_9',
  E_9_SPECIAL_SENIOR_ENLISTED: 'E_9_SPECIAL_SENIOR_ENLISTED',
  O_1: 'O_1_ACADEMY_GRADUATE',
  O_1_ACADEMY_GRADUATE: 'O_1_ACADEMY_GRADUATE',
  O_2: 'O_2',
  O_3: 'O_3',
  O_4: 'O_4',
  O_5: 'O_5',
  O_6: 'O_6',
  O_7: 'O_7',
  O_8: 'O_8',
  O_9: 'O_9',
  O_10: 'O_10',
  W_1: 'W_1',
  W_2: 'W_2',
  W_3: 'W_3',
  W_4: 'W_4',
  W_5: 'W_5',
  AVIATION_CADET: 'AVIATION_CADET',
  CIVILIAN_EMPLOYEE: 'CIVILIAN_EMPLOYEE',
  ACADEMY_CADET: 'ACADEMY_CADET',
  MIDSHIPMAN: 'MIDSHIPMAN',
};

export const ORDERS_PAY_GRADE_OPTIONS = {
  E_1: 'E-1',
  E_2: 'E-2',
  E_3: 'E-3',
  E_4: 'E-4',
  E_5: 'E-5',
  E_6: 'E-6',
  E_7: 'E-7',
  E_8: 'E-8',
  E_9: 'E-9',
  E_9_SPECIAL_SENIOR_ENLISTED: 'E-9 (Special Senior Enlisted)',
  O_1: 'O-1 or Service Academy Graduate',
  O_1_ACADEMY_GRADUATE: 'O-1 or Service Academy Graduate',
  O_2: 'O-2',
  O_3: 'O-3',
  O_4: 'O-4',
  O_5: 'O-5',
  O_6: 'O-6',
  O_7: 'O-7',
  O_8: 'O-8',
  O_9: 'O-9',
  O_10: 'O-10',
  W_1: 'W-1',
  W_2: 'W-2',
  W_3: 'W-3',
  W_4: 'W-4',
  W_5: 'W-5',
  AVIATION_CADET: 'Aviation Cadet',
  CIVILIAN_EMPLOYEE: 'Civilian Employee',
  ACADEMY_CADET: 'Service Academy Cadet',
  MIDSHIPMAN: 'Midshipman',
};

export const ORDERS_BRANCH_OPTIONS = {
  ARMY: 'Army',
  NAVY: 'Navy',
  MARINES: 'Marine Corps',
  AIR_FORCE: 'Air Force',
  COAST_GUARD: 'Coast Guard',
  SPACE_FORCE: 'Space Force',
  OTHER: 'Other',
};

export const ORDERS_DEPARTMENT_INDICATOR = {
  NAVY_AND_MARINES: 'Navy and Marine Corps',
  ARMY: 'Army',
  ARMY_CORPS_OF_ENGINEERS: 'Army Corps of Engineers',
  AIR_AND_SPACE_FORCE: 'Air Force and Space Force',
  COAST_GUARD: 'Coast Guard',
  OFFICE_OF_SECRETARY_OF_DEFENSE: 'Office of the Secretary of Defense',
};

export const RANK_GRADE_ASSOCIATIONS = {
  OTHER: [
    {
      abbv_rank: 'CIV',
      rank: 'Civilian',
      grade: 'CIVILIAN_EMPLOYEE',
    },
  ],
  NAVY: [
    {
      abbv_rank: 'ADM',
      rank: 'Admiral',
      grade: 'O_10',
    },
    {
      abbv_rank: 'VAD',
      rank: 'Vice Admiral',
      grade: 'O_9',
    },
    {
      abbv_rank: 'RADU',
      rank: 'Rear Admiral (Upper Half)',
      grade: 'O_8',
    },
    {
      abbv_rank: 'RADL',
      rank: 'Rear Admiral (Lower Half)',
      grade: 'O_7',
    },
    {
      abbv_rank: 'CPN',
      rank: 'Captain',
      grade: 'O_6',
    },
    {
      abbv_rank: 'CDR',
      rank: 'Commander',
      grade: 'O_5',
    },
    {
      abbv_rank: 'LCD',
      rank: 'Lieutenant Commander',
      grade: 'O_4',
    },
    {
      abbv_rank: 'LT',
      rank: 'Lieutenant',
      grade: 'O_3',
    },
    {
      abbv_rank: 'LTJG',
      rank: 'Lieutenant JG',
      grade: 'O_2',
    },
    {
      abbv_rank: 'ENS',
      rank: 'Ensign',
      grade: 'O_1_ACADEMY_GRADUATE',
    },
    {
      abbv_rank: 'MID',
      rank: 'Midshipman',
      grade: 'O_1_ACADEMY_GRADUATE',
    },
    {
      abbv_rank: 'WO5',
      rank: 'Chief Warrant Officer 5',
      grade: 'W_5',
    },
    {
      abbv_rank: 'WO4',
      rank: 'Chief Warrant Officer 4',
      grade: 'W_4',
    },
    {
      abbv_rank: 'WO3',
      rank: 'Chief Warrant Officer 3',
      grade: 'W_3',
    },
    {
      abbv_rank: 'WO2',
      rank: 'Chief Warrant Officer 2',
      grade: 'W_2',
    },
    {
      abbv_rank: 'WO1',
      rank: 'Warrant Officer 1',
      grade: 'W_1',
    },
    {
      abbv_rank: 'CPM',
      rank: 'Master Chief Petty Officer',
      grade: 'E_9',
    },
    {
      abbv_rank: 'MCPN',
      rank: 'Master Chief Petty Officer of the Navy',
      grade: 'E_9_SPECIAL_SENIOR_ENLISTED',
    },
    {
      abbv_rank: 'CPS',
      rank: 'Senior Chief Petty Officer',
      grade: 'E_8',
    },
    {
      abbv_rank: 'CPO',
      rank: 'Chief Petty Officer',
      grade: 'E_7',
    },
    {
      abbv_rank: 'PO1',
      rank: 'Petty Officer First Class',
      grade: 'E_6',
    },
    {
      abbv_rank: 'PO2',
      rank: 'Petty Officer Second Class',
      grade: 'E_5',
    },
    {
      abbv_rank: 'PO3',
      rank: 'Petty Officer Third Class',
      grade: 'E_4',
    },
    {
      abbv_rank: 'SN',
      rank: 'Seaman',
      grade: 'E_3',
    },
    {
      abbv_rank: 'SA',
      rank: 'Seaman Apprentice',
      grade: 'E_2',
    },
    {
      abbv_rank: 'SR',
      rank: 'Seaman Recruit',
      grade: 'E_1',
    },
    {
      abbv_rank: 'CIV',
      rank: 'Civilian',
      grade: 'CIVILIAN_EMPLOYEE',
    },
  ],
  SPACE_FORCE: [
    {
      abbv_rank: 'GEN',
      rank: 'General',
      grade: 'O_10',
    },
    {
      abbv_rank: 'LTG',
      rank: 'Lieutenant General',
      grade: 'O_9',
    },
    {
      abbv_rank: 'MG',
      rank: 'Major General',
      grade: 'O_8',
    },
    {
      abbv_rank: 'BG',
      rank: 'Brigadier General',
      grade: 'O_7',
    },
    {
      abbv_rank: 'COL',
      rank: 'Colonel',
      grade: 'O_6',
    },
    {
      abbv_rank: 'LTC',
      rank: 'Lieutenant Colonel',
      grade: 'O_5',
    },
    {
      abbv_rank: 'MAJ',
      rank: 'Major',
      grade: 'O_4',
    },
    {
      abbv_rank: 'CPT',
      rank: 'Captain',
      grade: 'O_3',
    },
    {
      abbv_rank: '1LT',
      rank: 'First Lieutenant',
      grade: 'O_2',
    },
    {
      abbv_rank: '2LT',
      rank: 'Second Lieutenant',
      grade: 'O_1_ACADEMY_GRADUATE',
    },
    {
      abbv_rank: 'CMA',
      rank: 'Chief Master Sergeant of the Space Force',
      grade: 'E_9_SPECIAL_SENIOR_ENLISTED',
    },
    {
      abbv_rank: 'CMS',
      rank: 'Chief Master Sergeant',
      grade: 'E_9',
    },
    {
      abbv_rank: 'SMS',
      rank: 'Senior Master Sergeant',
      grade: 'E_8',
    },
    {
      abbv_rank: 'MSG',
      rank: 'Master Sergeant',
      grade: 'E_7',
    },
    {
      abbv_rank: 'TSG',
      rank: 'Technical Sergeant',
      grade: 'E_6',
    },
    {
      abbv_rank: 'SGT',
      rank: 'Sergeant',
      grade: 'E_5',
    },
    {
      abbv_rank: 'SP4',
      rank: 'Specialist 4',
      grade: 'E_4',
    },
    {
      abbv_rank: 'SP3',
      rank: 'Specialist 3',
      grade: 'E_3',
    },
    {
      abbv_rank: 'SP2',
      rank: 'Specialist 2',
      grade: 'E_2',
    },
    {
      abbv_rank: 'SP1',
      rank: 'Specialist 1',
      grade: 'E_1',
    },
    {
      abbv_rank: 'CIV',
      rank: 'Civilian',
      grade: 'CIVILIAN_EMPLOYEE',
    },
  ],
  COAST_GUARD: [
    {
      abbv_rank: 'ADM',
      rank: 'Admiral',
      grade: 'O_10',
    },
    {
      abbv_rank: 'VAD',
      rank: 'Vice Admiral',
      grade: 'O_9',
    },
    {
      abbv_rank: 'RADU',
      rank: 'Rear Admiral (Upper Half)',
      grade: 'O_8',
    },
    {
      abbv_rank: 'RADL',
      rank: 'Rear Admiral (Lower Half)',
      grade: 'O_7',
    },
    {
      abbv_rank: 'CAPT',
      rank: 'Captain',
      grade: 'O_6',
    },
    {
      abbv_rank: 'CDR',
      rank: 'Commander',
      grade: 'O_5',
    },
    {
      abbv_rank: 'LCD',
      rank: 'Lieutenant Commander',
      grade: 'O_4',
    },
    {
      abbv_rank: 'LT',
      rank: 'Lieutenant',
      grade: 'O_3',
    },
    {
      abbv_rank: 'LTJG',
      rank: 'Lieutenant JG',
      grade: 'O_2',
    },
    {
      abbv_rank: 'MID',
      rank: 'Midshipman',
      grade: 'O_1_ACADEMY_GRADUATE',
    },
    {
      abbv_rank: 'ENS',
      rank: 'Ensign',
      grade: 'O_1_ACADEMY_GRADUATE',
    },
    {
      abbv_rank: 'CWO4',
      rank: 'Chief Warrant Officer 4',
      grade: 'W_4',
    },
    {
      abbv_rank: 'CWO3',
      rank: 'Chief Warrant Officer 3',
      grade: 'W_3',
    },
    {
      abbv_rank: 'CWO2',
      rank: 'Chief Warrant Officer 2',
      grade: 'W_2',
    },
    {
      abbv_rank: 'CPM',
      rank: 'Master Chief Petty Officer',
      grade: 'E_9',
    },
    {
      abbv_rank: 'MCPG',
      rank: 'Master Chief Petty Officer of the Coast Guard',
      grade: 'E_9_SPECIAL_SENIOR_ENLISTED',
    },
    {
      abbv_rank: 'CPS',
      rank: 'Senior Chief Petty Officer',
      grade: 'E_8',
    },
    {
      abbv_rank: 'CPO',
      rank: 'Chief Petty Officer',
      grade: 'E_7',
    },
    {
      abbv_rank: 'PO1',
      rank: 'Petty Officer First Class',
      grade: 'E_6',
    },
    {
      abbv_rank: 'PO2',
      rank: 'Petty Officer Second Class',
      grade: 'E_5',
    },
    {
      abbv_rank: 'PO3',
      rank: 'Petty Officer Third Class',
      grade: 'E_4',
    },
    {
      abbv_rank: 'SN',
      rank: 'Seaman',
      grade: 'E_3',
    },
    {
      abbv_rank: 'SA',
      rank: 'Seaman Apprentice',
      grade: 'E_2',
    },
    {
      abbv_rank: 'SR',
      rank: 'Seaman Recruit',
      grade: 'E_1',
    },
    {
      abbv_rank: 'CIV',
      rank: 'Civilian',
      grade: 'CIVILIAN_EMPLOYEE',
    },
  ],
  MARINES: [
    {
      abbv_rank: 'GEN',
      rank: 'General',
      grade: 'O_10',
    },
    {
      abbv_rank: 'LTG',
      rank: 'Lieutenant General',
      grade: 'O_9',
    },
    {
      abbv_rank: 'MG',
      rank: 'Major General',
      grade: 'O_8',
    },
    {
      abbv_rank: 'BG',
      rank: 'Brigadier General',
      grade: 'O_7',
    },
    {
      abbv_rank: 'COL',
      rank: 'Colonel',
      grade: 'O_6',
    },
    {
      abbv_rank: 'LTC',
      rank: 'Lieutenant Colonel',
      grade: 'O_5',
    },
    {
      abbv_rank: 'MAJ',
      rank: 'Major',
      grade: 'O_4',
    },
    {
      abbv_rank: 'CPT',
      rank: 'Captain',
      grade: 'O_3',
    },
    {
      abbv_rank: '1LT',
      rank: 'First Lieutenant',
      grade: 'O_2',
    },
    {
      abbv_rank: '2LT',
      rank: 'Second Lieutenant',
      grade: 'O_1_ACADEMY_GRADUATE',
    },
    {
      abbv_rank: 'CW5',
      rank: 'Chief Warrant Officer 5',
      grade: 'W_5',
    },
    {
      abbv_rank: 'CW4',
      rank: 'Chief Warrant Officer 4',
      grade: 'W_4',
    },
    {
      abbv_rank: 'CW3',
      rank: 'Chief Warrant Officer 3',
      grade: 'W_3',
    },
    {
      abbv_rank: 'CW2',
      rank: 'Chief Warrant Officer 2',
      grade: 'W_2',
    },
    {
      abbv_rank: 'WO1',
      rank: 'Warrant Officer 1',
      grade: 'W_1',
    },
    {
      abbv_rank: 'SMG',
      rank: 'Sergeant Major',
      grade: 'E_9',
    },
    {
      abbv_rank: 'SMM',
      rank: 'Sergeant Major of the Marine Corps',
      grade: 'E_9_SPECIAL_SENIOR_ENLISTED',
    },
    {
      abbv_rank: 'MGS',
      rank: 'Master Gunnery Sergeant',
      grade: 'E_9',
    },
    {
      abbv_rank: 'MSG',
      rank: 'Master Sergeant',
      grade: 'E_8',
    },
    {
      abbv_rank: '1ST',
      rank: '1st Sergeant',
      grade: 'E_8',
    },
    {
      abbv_rank: 'GYS',
      rank: 'Gunnery Sergeant',
      grade: 'E_7',
    },
    {
      abbv_rank: 'SSG',
      rank: 'Staff Sergeant',
      grade: 'E_6',
    },
    {
      abbv_rank: 'SGT',
      rank: 'Sergeant',
      grade: 'E_5',
    },
    {
      abbv_rank: 'CPL',
      rank: 'Corporal',
      grade: 'E_4',
    },
    {
      abbv_rank: 'LCP',
      rank: 'Lance Corporal',
      grade: 'E_3',
    },
    {
      abbv_rank: 'PFC',
      rank: 'Private First Class',
      grade: 'E_2',
    },
    {
      abbv_rank: 'PVT',
      rank: 'Private',
      grade: 'E_1',
    },
    {
      abbv_rank: 'CIV',
      rank: 'Civilian',
      grade: 'CIVILIAN_EMPLOYEE',
    },
  ],
  AIR_FORCE: [
    {
      abbv_rank: 'GEN',
      rank: 'General',
      grade: 'O_10',
    },
    {
      abbv_rank: 'LTG',
      rank: 'Lieutenant General',
      grade: 'O_9',
    },
    {
      abbv_rank: 'MG',
      rank: 'Major General',
      grade: 'O_8',
    },
    {
      abbv_rank: 'BG',
      rank: 'Brigadier General',
      grade: 'O_7',
    },
    {
      abbv_rank: 'COL',
      rank: 'Colonel',
      grade: 'O_6',
    },
    {
      abbv_rank: 'LTC',
      rank: 'Lieutenant Colonel',
      grade: 'O_5',
    },
    {
      abbv_rank: 'MAJ',
      rank: 'Major',
      grade: 'O_4',
    },
    {
      abbv_rank: 'CPT',
      rank: 'Captain',
      grade: 'O_3',
    },
    {
      abbv_rank: '1LT',
      rank: 'First Lieutenant',
      grade: 'O_2',
    },
    {
      abbv_rank: 'AVC',
      rank: 'Aviation Cadet',
      grade: 'O_1_ACADEMY_GRADUATE',
    },
    {
      abbv_rank: '2LT',
      rank: 'Second Lieutenant',
      grade: 'O_1_ACADEMY_GRADUATE',
    },
    {
      abbv_rank: 'AFC',
      rank: 'Air Force Academy Cadet',
      grade: 'O_1_ACADEMY_GRADUATE',
    },
    {
      abbv_rank: 'CMS',
      rank: 'Chief Master Sergeant',
      grade: 'E_9',
    },
    {
      abbv_rank: 'SMS',
      rank: 'Senior Master Sergeant',
      grade: 'E_8',
    },
    {
      abbv_rank: 'MSG',
      rank: 'Master Sergeant',
      grade: 'E_7',
    },
    {
      abbv_rank: 'TSG',
      rank: 'Technical Sergeant',
      grade: 'E_6',
    },
    {
      abbv_rank: 'SSG',
      rank: 'Staff Sergeant',
      grade: 'E_5',
    },
    {
      abbv_rank: 'SRA',
      rank: 'Senior Airman',
      grade: 'E_4',
    },
    {
      abbv_rank: 'A1C',
      rank: 'Airman First Class',
      grade: 'E_3',
    },
    {
      abbv_rank: 'AMN',
      rank: 'Airman',
      grade: 'E_2',
    },
    {
      abbv_rank: 'AB',
      rank: 'Airman Basic',
      grade: 'E_1',
    },
    {
      abbv_rank: 'CIV',
      rank: 'Civilian',
      grade: 'CIVILIAN_EMPLOYEE',
    },
  ],
  ARMY: [
    {
      abbv_rank: 'GEN',
      rank: 'General',
      grade: 'O_10',
    },
    {
      abbv_rank: 'LTG',
      rank: 'Lieutenant General',
      grade: 'O_9',
    },
    {
      abbv_rank: 'MG',
      rank: 'Major General',
      grade: 'O_8',
    },
    {
      abbv_rank: 'BG',
      rank: 'Brigadier General',
      grade: 'O_7',
    },
    {
      abbv_rank: 'COL',
      rank: 'Colonel',
      grade: 'O_6',
    },
    {
      abbv_rank: 'LTC',
      rank: 'Lieutenant Colonel',
      grade: 'O_5',
    },
    {
      abbv_rank: 'MAJ',
      rank: 'Major',
      grade: 'O_4',
    },
    {
      abbv_rank: 'CPT',
      rank: 'Captain',
      grade: 'O_3',
    },
    {
      abbv_rank: '1LT',
      rank: 'First Lieutenant',
      grade: 'O_2',
    },
    {
      abbv_rank: '2LT',
      rank: 'Second Lieutenant',
      grade: 'O_1_ACADEMY_GRADUATE',
    },
    {
      abbv_rank: 'OC',
      rank: 'Officer Candidate',
      grade: 'O_1_ACADEMY_GRADUATE',
    },
    {
      abbv_rank: 'CDT',
      rank: 'Cadet',
      grade: 'O_1_ACADEMY_GRADUATE',
    },
    {
      abbv_rank: 'CW5',
      rank: 'Chief Warrant Officer 5',
      grade: 'W_5',
    },
    {
      abbv_rank: 'CW4',
      rank: 'Chief Warrant Officer 4',
      grade: 'W_4',
    },
    {
      abbv_rank: 'CW3',
      rank: 'Chief Warrant Officer 3',
      grade: 'W_3',
    },
    {
      abbv_rank: 'CW2',
      rank: 'Chief Warrant Officer 2',
      grade: 'W_2',
    },
    {
      abbv_rank: 'WO1',
      rank: 'Warrant Officer 1',
      grade: 'W_1',
    },
    {
      abbv_rank: 'SMA',
      rank: 'Sergeant Major of the Army',
      grade: 'E_9_SPECIAL_SENIOR_ENLISTED',
    },
    {
      abbv_rank: 'SGM',
      rank: 'Sergeant Major',
      grade: 'E_9',
    },
    {
      abbv_rank: 'CSM',
      rank: 'Command Sergeant Major',
      grade: 'E_9',
    },
    {
      abbv_rank: 'MSG',
      rank: 'Master Sergeant',
      grade: 'E_8',
    },
    {
      abbv_rank: '1ST',
      rank: '1st Sergeant',
      grade: 'E_8',
    },
    {
      abbv_rank: 'PSG',
      rank: 'Platoon Sergeant',
      grade: 'E_7',
    },
    {
      abbv_rank: 'SFC',
      rank: 'Sergeant First Class',
      grade: 'E_7',
    },
    {
      abbv_rank: 'SSG',
      rank: 'Staff Sergeant',
      grade: 'E_6',
    },
    {
      abbv_rank: 'SGT',
      rank: 'Sergeant',
      grade: 'E_5',
    },
    {
      abbv_rank: 'CPL',
      rank: 'Corporal',
      grade: 'E_4',
    },
    {
      abbv_rank: 'SPC',
      rank: 'Specialist',
      grade: 'E_4',
    },
    {
      abbv_rank: 'PFC',
      rank: 'Private First Class',
      grade: 'E_3',
    },
    {
      abbv_rank: 'PV2',
      rank: 'Private',
      grade: 'E_2',
    },
    {
      abbv_rank: 'PV1',
      rank: 'Private',
      grade: 'E_1',
    },
    {
      abbv_rank: 'PVT',
      rank: 'Private',
      grade: 'E_1',
    },
    {
      abbv_rank: 'CIV',
      rank: 'Civilian',
      grade: 'CIVILIAN_EMPLOYEE',
    },
  ],
};

export const formatRankGradeDisplayValue = ({ rank, grade }) => {
  return `${rank} / ${grade}`;
};

export const rankOptionValuesByAffiliation = (affiliation) => {
  const affiliationPaygradeRankEntries = Object.fromEntries(
    (RANK_GRADE_ASSOCIATIONS[affiliation] ?? []).map((e) => [
      e.abbv_rank,
      { value: formatRankGradeDisplayValue({ rank: e.abbv_rank, grade: ORDERS_PAY_GRADE_OPTIONS[e.grade] }), ...e },
    ]),
  );
  return affiliationPaygradeRankEntries;
};
