import React from 'react';
import { Tag } from '@trussworks/react-uswds';

import LeftNav from './index';

import { faPlus } from '@fortawesome/free-solid-svg-icons';

// Left Nav
export default {
  title: 'Components|Left Nav',
  component: LeftNav,
  parameters: {
    abstract: {
      url: 'https://share.goabstract.com/6e8668b7-5562-4894-a661-648ab4883d8f?mode=design',
    },
  },
};

export const Basic = () => (
  <div id="l-nav" style={{ padding: '20px', background: '#f0f0f0' }}>
    <LeftNav>
      <a href="#" className="active">
        Requested Shipments
        <Tag className="usa-tag--alert usa-tag--alert--small">
          <FontAwesomeIcon icon={faPlus} />
        </Tag>
      </a>
      <a href="#orders-anchor">
        Orders
        <Tag className="usa-tag--teal">INTL</Tag>
      </a>
      <a href="#">Allowances</a>
      <a href="#">
        Customer Info
        <Tag>3</Tag>
      </a>
    </LeftNav>
  </div>
);

export const Layout = () => (
  <div style={{ padding: '20px', background: '#f0f0f0', display: 'flex' }}>
    <div style={{ position: 'sticky', alignSelf: 'flex-start', marginTop: '26px', top: '26px' }}>
      <LeftNav>
        <a href="#" className="active">
          Requested Shipments
          <Tag className="usa-tag--alert usa-tag--alert--small">
            <FontAwesomeIcon icon={faPlus} />
          </Tag>
        </a>
        <a href="#orders-anchor">
          Orders
          <Tag className="usa-tag--teal">INTL</Tag>
        </a>
        <a href="#">Allowances</a>
        <a href="#">
          Customer Info
          <Tag>3</Tag>
        </a>
      </LeftNav>
    </div>
    <div
      className="content container container--accent--hhg"
      style={{ position: 'relative', width: '85vw', marginTop: '0' }}
    >
      <ul>
        <h2>Fixed positioning behavior of the tertiary nav</h2>
        <li>
          Sections of the page are anchored, selecting an item on the tab&nbsp;
          <a
            href="https://xc9rwh.axshare.com/#id=gro6ti&p=move_details_dir_a&g=1"
            target="_blank"
            rel="noopener noreferrer"
          >
            will take you to the respective page.
          </a>
        </li>
        <li>
          The tertiary nav is pinned to the browser and scrolls with the page using &nbsp;
          <code>position: sticky;</code>
          &nbsp;. This is not supported in IE and will need to be polyfilled.
        </li>
        <li>
          click the&nbsp;
          <b>Orders</b>
          &nbsp;tab in the left nav to demo it. In Storybook, you have to have to&nbsp;
          <b>Open canvas in a new tab</b>
          &nbsp;for it to work.
        </li>
        <li>
          in this example,&nbsp;
          <code>sidebar</code>
          &nbsp;and the rest of the content are being laid out using&nbsp;
          <a href="https://css-tricks.com/snippets/css/a-guide-to-flexbox/" target="_blank" rel="noopener noreferrer">
            flexbox
          </a>
          .&nbsp;
          <code>sidebar</code>
          &nbsp;has a&nbsp;
          <code>width</code>
          &nbsp;of&nbsp;
          <b>230px</b>
          &nbsp;set by default. If you wrap &nbsp;
          <code>sidebar</code>
          &nbsp; and your intended right side of the page content in a div and apply
          <code>display: flex;</code>
          &nbsp;and set your content&apos;s wrapping div to&nbsp;
          <b>85vw</b>
          &nbsp;, you should get the intended layout.
        </li>
        <li>
          <b>Smooth Scrolling</b>
          is currently implemented in CSS by applying&nbsp;
          <code>scroll-behavior: smooth;</code>
          &nbsp;to the&nbsp;
          <code>html</code>
          &nbsp;element.&nbsp;
          <b>This does not have Safari support</b>
          &nbsp;and we may need to explore js solutions at least as a fallback.
        </li>
      </ul>
      <div
        style={{
          background: 'lightgray',
          display: 'flex',
          flexDirection: 'column',
          justifyContent: 'center',
          alignItems: 'center',
          height: '768px',
        }}
      >
        <p>
          hi, I&apos;m &nbsp;
          <b>scrollable content</b>
          .&nbsp;
        </p>
        <p id="orders-anchor">I&apos;m an anchor point.</p>
      </div>
    </div>
  </div>
);
