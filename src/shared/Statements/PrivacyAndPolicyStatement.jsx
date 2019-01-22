import React from 'react';
import './statements.css';

function PrivacyPolicy() {
  return (
    <div className="usa-grid">
      <div className="usa-width-two-thirds statement-content">
        <h1>Privacy & Security Policy</h1>
        <p>
          Your privacy is important to us. By using{' '}
          <a href="http://my.move.mil/" target="_blank" rel="noopener noreferrer">
            my.move.mil
          </a>
          , you authorize us to share your data with other government entities to facilitate your personal property
          relocation.
        </p>
        <p>
          {' '}
          All records are stored electronically in a database in DoD’s{' '}
          <a href="https://aws.amazon.com/" target="_blank" rel="noopener noreferrer">
            Amazon Web Services (AWS) environment
          </a>
          . All records are encrypted while being stored (at rest) and when the data is transferred between systems (in
          transit).
        </p>
        <p>
          <a href="https://my.move.mil/" target="_blank" rel="noopener noreferrer">
            My.move.mil
          </a>{' '}
          uses industry best practices in information security, drawing upon DoD, HIPAA, and PCI standards for
          information assurance.
        </p>
        <p>
          The information you provide to access your{' '}
          <a href="https://my.move.mil/" target="_blank" rel="noopener noreferrer">
            my.move.mil
          </a>{' '}
          account is collected pursuant to 6 USC § 1523 (b)(1)(A)-(E), the{' '}
          <a
            href="https://www.gpo.gov/fdsys/pkg/PLAW-107publ347/html/PLAW-107publ347.htm"
            target="_blank"
            rel="noopener noreferrer"
          >
            E-Government Act of 2002 (44 USC § 3501)
          </a>
          , and 40 USC § 501.
        </p>
        <p>
          Access and accessibility is important to us.{' '}
          <a href="https://my.move.mil/" target="_blank" rel="noopener noreferrer">
            My.move.mil
          </a>{' '}
          is built to be accessible and compliant with Section 508. If you discover an issue with this site’s
          accessibility, please report it to your local Personal Property Office (PPPO).
        </p>
      </div>
    </div>
  );
}

export default PrivacyPolicy;
