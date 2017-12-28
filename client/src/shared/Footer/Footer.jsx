import React from 'react';

function Footer() {
  return (
    <footer class="usa-footer usa-footer-medium" role="contentinfo">
      <div class="usa-grid usa-footer-return-to-top">
        <a href="">Return to top</a>
      </div>
      <div class="usa-footer-primary-section">
        <div class="usa-grid-full">
          <nav class="usa-footer-nav">
            <ul class="usa-unstyled-list">
              <li class="usa-width-one-fourth usa-footer-primary-content">
                <a class="usa-footer-primary-link" href="https://www.move.mil/">
                  Move.mil
                </a>
              </li>
              <li class="usa-width-one-fourth usa-footer-primary-content">
                <a class="usa-footer-primary-link" href="">
                  Help Me
                </a>
              </li>
              <li class="usa-width-one-fourth usa-footer-primary-content">
                <a class="usa-footer-primary-link" href="">
                  Site policies (example)
                </a>
              </li>
            </ul>
          </nav>
        </div>
      </div>

      <div class="usa-footer-secondary_section">
        <div class="usa-grid">
          <div class="usa-footer-logo usa-width-one-half">
            <a href="https://www.ustranscom.mil/">
              <img
                class="usa-footer-logo-img"
                src=""
                alt="United States Transportation Command Emblem"
              />
              <br />
              <h3 class="usa-footer-big-logo-heading">USTRANSCOM</h3>
              <span>United States Transportation Command</span>
              <br />
            </a>
          </div>
          <div class="usa-footer-contact-links usa-width-one-half">
            <a class="usa-link-twitter" href="">
              <span>Twitter</span>
            </a>
            <a class="usa-link-rss" href="">
              <span>RSS</span>
            </a>
            <address>
              <h3 class="usa-footer-contact-heading">Contact Us</h3>
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
