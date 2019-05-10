import React from 'react';
import transcomEmblem from 'shared/images/transcom-emblem.svg';
import { Link } from 'react-router-dom';

function Footer() {
  return (
    <footer className="usa-footer usa-footer-medium" role="contentinfo">
      <div className="usa-grid usa-footer-return-to-top" />
      <div className="usa-footer-primary-section">
        <div className="usa-grid-full">
          <nav className="usa-footer-nav">
            <ul className="usa-unstyled-list">
              <li className="usa-width-one-fourth usa-footer-primary-content">
                <a className="usa-footer-primary-link" href="https://www.move.mil/">
                  Move.mil
                </a>
              </li>
              <li className="usa-width-one-fourth usa-footer-primary-content">
                <a className="usa-footer-primary-link" href="mailto:transcom.scott.tcj5j4.mbx.ppcf@mail.mil">
                  Help Me
                </a>
              </li>
              <li className="usa-width-one-fourth usa-footer-primary-content">
                <Link className="usa-footer-primary-link" to="/accessibility">
                  Accessibility
                </Link>
              </li>
              <li className="usa-width-one-fourth usa-footer-primary-content">
                <Link className="usa-footer-primary-link" to="/privacy-and-security-policy">
                  Privacy and Security Policy
                </Link>
              </li>
            </ul>
          </nav>
        </div>
      </div>
      <div className="usa-footer-secondary_section">
        <div className="usa-grid">
          <div className="usa-footer-logo usa-width-one-half">
            <a href="https://www.ustranscom.mil/">
              <img
                className="usa-footer-logo-img _fix-ie-11-height"
                src={transcomEmblem}
                alt="United States Transportation Command Emblem"
              />
              <br />
              <h3 className="usa-footer-big-logo-heading">USTRANSCOM</h3>
              <span>United States Transportation Command</span>
              <br />
            </a>
          </div>
          <div className="usa-footer-contact-links usa-width-one-half">
            <a
              className="usa-link-twitter"
              href="https://twitter.com/us_transcom"
              target="_blank"
              rel="noopener noreferrer"
            >
              <span>Twitter</span>
            </a>
            <a
              className="usa-link-facebook"
              href="https://www.facebook.com/USTRANSCOM/"
              target="_blank"
              rel="noopener noreferrer"
            >
              <span>Facebook</span>
            </a>
            <address>
              <h3 className="usa-footer-contact-heading">Contact Us</h3>
              <p>
                <a href="https://move.mil/customer-service">Customer service</a>
              </p>
            </address>
          </div>
        </div>
      </div>
    </footer>
  );
}

export default Footer;
