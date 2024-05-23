import React from 'react';

import { useTitle } from 'hooks/custom';
import './statements.css';

function PrivacyPolicy() {
  useTitle('Privacy & Security Policy');
  return (
    <div className="usa-grid">
      <div className="usa-width-two-thirds statement-content">
        <h1>Privacy & Security Policy</h1>
        <p>
          1. This system is covered by Systems of Record Notice FRTRANSCOM 01 DoD - Defense Transportation System (DTS)
          Records dated(September 26, 2014, 79 FR 57893).
        </p>
        <p>
          2. AUTHORITY: Public Law 100-562, Imported Vehicle Safety Compliance Act of 1988; 5 U.S.C. 5726, Storage
          Expenses, Household Goods and Personal Effects; 10 U.S.C. 113, Secretary of Defense; 10 U.S.C. 3013, Secretary
          of the Army; 10 U.S.C. 5013, Secretary of the Navy; 10 U.S.C. 8013, Secretary of the Air Force; 19 U.S.C.
          1498, Entry Under Regulations; 37 U.S.C. 406, Travel and Transportation Allowances, Dependents, Baggage and
          Household Effects; Federal Acquisition Regulation (FAR); Joint Federal Travel Regulation (JTR), Volumes I and
          II; DoD Directive 4500.9E, Transportation and Traffic Management; DoD Directive 5158.4, United States
          Transportation Command; DoD Instruction 4500.42, DoD Transportation Reservation and Ticketing Services; DoD
          Regulation 4140.1, DoD Materiel Management Regulation; DoD Regulation 4500.9, Defense Transportation
          Regulation; and DoD Regulation 4515.13-R, Air Transportation Eligibility. DPS - Executive Order 10450, 9397;
          and Public Law 99-474, the Computer Fraud and Abuse Act. Privacy Act Information -The information accessed
          through this system is For Official Use Only and must be protected IAW DoD Directive 5400.11 and DoD
          5400.11-R, DoD Privacy Program. Authority: The Privacy Act of 1974, as amended. 5 U.S.C. 552a.
        </p>
        <p>
          PURPOSE: PII is collected to effectuate movement of household goods and personal property shipments, to
          include payment of carrier and service-members.
        </p>
        <p>
          ROUTINE USES: All collected PII will be disclosed to agencies and activities of the Department of Defense for
          the purposes of arranging and facilitating tasks directly related to effecting movement of household goods and
          personal property to include facilitating payment to Service-members. Certain information, limited to rolodex
          contact information, is disclosed to contracted Transportation Services Providers in order to effect the
          movement of household goods and personal property. Use of information in this system is restricted to MilMove
          account holders and disclosure is prohibited without the written consent of the Defense Personal Property
          Management Office (DPMO).
        </p>
        <p>
          DISCLOSURE: Voluntary. However, failure to provide the requested information may result in the Service-member
          being unable to use the MilMove system to effectuate their household goods or personal property movement.
        </p>
      </div>
    </div>
  );
}

export default PrivacyPolicy;
