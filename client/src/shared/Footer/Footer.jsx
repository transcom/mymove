import React from 'react';
import transcomEmblem from '../../transcom-emblem.svg';

function Footer() {
  return (
    <footer className="usa-footer usa-footer-medium" role="contentinfo">
      <div className="usa-grid usa-footer-return-to-top">
        <a href>Return to top</a>
      </div>
      <div className="usa-footer-primary-section">
        <div className="usa-grid-full">
          <nav className="usa-footer-nav">
            <ul className="usa-unstyled-list">
              <li className="usa-width-one-fourth usa-footer-primary-content">
                <a
                  className="usa-footer-primary-link"
                  href="https://www.move.mil/"
                >
                  Move.mil
                </a>
              </li>
              <li className="usa-width-one-fourth usa-footer-primary-content">
                <a className="usa-footer-primary-link" href>
                  Help Me
                </a>
              </li>
              <li className="usa-width-one-fourth usa-footer-primary-content">
                <a className="usa-footer-primary-link" href>
                  Site policies (example)
                </a>
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
                className="usa-footer-logo-img"
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
            <a className="usa-link-twitter" href>
              <span>Twitter</span>
            </a>
            <a className="usa-link-rss" href>
              <span>RSS</span>
            </a>
            <address>
              <h3 className="usa-footer-contact-heading">Contact Us</h3>
              <p>(800) CALL-GOVT</p>
              <a href="mailto:info@agency.gov">info@agency.gov</a>
            </address>
          </div>
        </div>
      </div>
    </footer>
  );
}

export default Footer;
