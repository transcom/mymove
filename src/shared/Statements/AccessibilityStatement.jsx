import React from 'react';
import './statements.css';

function PrivacyPolicy() {
  return (
    <div className="usa-grid">
      <div className="usa-width-two-thirds statement-content">
        <h1>508 Compliance</h1>
        <p>
          If your issue involves log in access, password recovery, or other technical issues, contact the administrator
          for the website in question, please report it to your local Personal Property Office (PPPO).
        </p>
        <p>
          The U.S. Department of Defense is committed to making its electronic and information technologies accessible
          to individuals with disabilities in accordance with{' '}
          <a
            href="https://www.access-board.gov/the-board/laws/rehabilitation-act-of-1973#508%20"
            target="_blank"
            rel="noopener noreferrer"
          >
            Section 508 of the Rehabilitation Act (29 U.S.C. 794d)
          </a>
          , as amended in 1998.
        </p>
        <p>
          For persons with disabilities experiencing difficulties accessing content on a particular website, please use
          the{' '}
          <a
            href="http://dodcio.defense.gov/DoDSection508/Section508Form.aspx"
            target="_blank"
            rel="noopener noreferrer"
          >
            DoD Section 508 Form
          </a>
          . In this form, please indicate the nature of your accessibility issue/problem and your contact information so
          we can address your issue or question.
        </p>
        <p>
          For more information about Section 508, please visit the{' '}
          <a href="http://dodcio.defense.gov/DoDSection508.aspx" target="_blank" rel="noopener noreferrer">
            {' '}
            DoD Section 508 website.
          </a>
        </p>
      </div>
    </div>
  );
}

export default PrivacyPolicy;
