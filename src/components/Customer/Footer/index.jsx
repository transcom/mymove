import React from 'react';
import { Link } from 'react-router-dom';
import iconFacebook from 'uswds/src/img/usa-icons/facebook.svg';
import iconTwitter from 'uswds/src/img/usa-icons/twitter.svg';

import transcomEmblem from 'shared/images/transcom-emblem.svg';

function Footer() {
  return (
    <footer className="usa-footer" role="contentinfo">
      <div className="grid-container usa-footer__return-to-top" />
      <div className="usa-footer__primary-section">
        <div className="usa-footer__primary-container grid-row">
          <div className="grid-col-12">
            <nav className="usa-footer__nav">
              <ul className="grid-row grid-gap">
                <li className="mobile-lg:grid-col-6 desktop:grid-col-auto usa-footer__primary-content">
                  <a className="usa-footer__primary-link" href="https://www.move.mil/">
                    Move.mil
                  </a>
                </li>
                <li className="mobile-lg:grid-col-6 desktop:grid-col-auto usa-footer__primary-content">
                  <a className="usa-footer__primary-link" href="mailto:transcom.scott.tcj5j4.mbx.ppcf@mail.mil">
                    Help Me
                  </a>
                </li>
                <li className="mobile-lg:grid-col-6 desktop:grid-col-auto usa-footer__primary-content">
                  <Link className="usa-footer__primary-link" to="/accessibility">
                    Accessibility
                  </Link>
                </li>
                <li className="mobile-lg:grid-col-6 desktop:grid-col-auto usa-footer__primary-content">
                  <Link className="usa-footer__primary-link" to="/privacy-and-security-policy">
                    Privacy and Security Policy
                  </Link>
                </li>
              </ul>
            </nav>
          </div>
        </div>
      </div>
      <div className="usa-footer__secondary-section">
        <div className="grid-container">
          <div className="grid-row grid-gap">
            <div className="usa-footer__logo grid-row mobile-lg:grid-col-6 mobile-lg:grid-gap-2">
              <div className="mobile-lg:grid-col-auto">
                <img
                  className="usa-footer__logo-img _fix-ie-11-height"
                  src={transcomEmblem}
                  alt="United States Transportation Command Emblem"
                />
              </div>
              <div className="mobile-lg:grid-col-auto">
                <a href="https://www.ustranscom.mil/">
                  <h3 className="usa-footer__logo-heading">USTRANSCOM</h3>
                  <span>United States Transportation Command</span>
                </a>
              </div>
            </div>
            <div className="usa-footer__contact-links mobile-lg:grid-col-6">
              <div className="usa-footer__social-links grid-row grid-gap-1">
                <div className="grid-col-auto">
                  <a
                    className="usa-social-link"
                    href="https://twitter.com/us_transcom"
                    target="_blank"
                    rel="noopener noreferrer"
                  >
                    <img className="usa-social-link__icon" src={iconTwitter} alt="Twitter" />
                  </a>
                </div>
                <div className="grid-col-auto">
                  <a
                    className="usa-social-link"
                    href="https://www.facebook.com/USTRANSCOM/"
                    target="_blank"
                    rel="noopener noreferrer"
                  >
                    <img className="usa-social-link__icon" src={iconFacebook} alt="Facebook" />
                  </a>
                </div>
              </div>
              <h3 className="usa-footer__contact-heading" data-testid="contact-footer">
                Contact Us
              </h3>
              <address className="usa-footer__address">
                <div className="usa-footer__contact-info grid-row grid-gap">
                  <a href="https://move.mil/customer-service">Customer service</a>
                </div>
              </address>
            </div>
          </div>
        </div>
      </div>
    </footer>
  );
}

export default Footer;
