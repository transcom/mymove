import React from 'react';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';

import { ReactComponent as ExternalLink } from 'shared/icon/external-link.svg';

// Icons
export default {
  title: 'Global/Icons',
  parameters: {
    abstract: {
      url: 'https://share.goabstract.com/eabef4a2-603e-4c3d-b249-0e580f1c8306?mode=design',
    },
  },
};

export const all = () => (
  <div style={{ padding: '20px', background: '#f0f0f0' }}>
    <h3>Icons</h3>
    <div
      id="icons"
      style={{
        display: 'grid',
        gridTemplateColumns: `repeat( auto-fit, minmax(150px, 1fr)`,
        gridTemplateRows: `repeat(5, 1fr)`,
      }}
    >
      <div>
        <FontAwesomeIcon icon="file" />
        <code>documents | icon=&quot;file&quot;</code>
      </div>
      <div>
        <FontAwesomeIcon icon="pen" />
        <code>edit | icon=&quot;pen&quot;</code>
      </div>
      <div>
        <FontAwesomeIcon icon="plus" />
        <code>add | icon=&quot;plus&quot;</code>
      </div>
      <div>
        <FontAwesomeIcon icon="chevron-left" />
        <code>chevron-left | icon=&quot;chevron-left&quot;</code>
      </div>
      <div>
        <FontAwesomeIcon icon="chevron-right" />
        <code>chevron-right | icon=&quot;chevron-right&quot;</code>
      </div>
      <div>
        <FontAwesomeIcon icon="chevron-down" />
        <code>chevron-down | icon=&quot;chevron-down&quot;</code>
      </div>
      <div>
        <FontAwesomeIcon icon="chevron-up" />
        <code>chevron-up | icon=&quot;chevron-up</code>
      </div>
      <div>
        <FontAwesomeIcon icon="check" />
        <code>checkmark | icon=&quot;check&quot;</code>
      </div>
      <div>
        <FontAwesomeIcon icon="times" />
        <code> x | icon=&quot;times&quot;</code>
      </div>
      <div>
        <FontAwesomeIcon icon={['far', 'circle-check']} />
        <code>accept (regular) | {`icon={['far', 'circle-check']}`}</code>
      </div>
      <div>
        <FontAwesomeIcon icon={['far', 'times-circle']} />
        <code>reject (regular) | {`icon={['far', 'times-circle']}`}</code>
      </div>
      <div>
        <FontAwesomeIcon icon="search-plus" />
        <code>zoom in | icon=&quot;search-plus&quot;</code>
      </div>
      <div>
        <FontAwesomeIcon icon="search-minus" />
        <code>zoom out | icon=&quot;search-minus&quot;</code>
      </div>
      <div>
        <FontAwesomeIcon icon="redo-alt" />
        <code>rotate clockwise | icon=&quot;redo-alt&quot;</code>
      </div>
      <div>
        <FontAwesomeIcon icon="undo-alt" />
        <code>rotate counter clockwise | icon=&quot;undo-alt&quot;</code>
      </div>
      <div>
        <FontAwesomeIcon icon="lock" />
        <code>lock | icon=&quot;lock&quot;</code>
      </div>
      <div>
        <FontAwesomeIcon icon="map-marker-alt" />
        <code>map pin | icon=&quot;map-marker-alt&quot;</code>
      </div>
      <div>
        <FontAwesomeIcon icon="arrow-right" />
        <code>arrow right | icon=&quot;arrow-right&quot;</code>
      </div>
      <div>
        <FontAwesomeIcon icon="arrow-left" />
        <code>arrow left| icon=&quot;arrow-left&quot;</code>
      </div>
      <div>
        <FontAwesomeIcon icon="th-list" />
        <code>doc menu | icon=&quot;th-list&quot;</code>
      </div>
      <div>
        <FontAwesomeIcon icon={['far', 'calendar']} />
        <code>calendar | {`icon={['far', 'calendar']}`}</code>
      </div>
      <div>
        <FontAwesomeIcon icon={['far', 'circle-question']} />
        <code>question circle (regular) | {`icon={['far', 'circle-question']}`}</code>
      </div>
      <div>
        <FontAwesomeIcon icon="circle-question" />
        <code>question circle (solid) | icon=&quot;circle-question&quot;</code>
      </div>
      <div>
        <FontAwesomeIcon icon="phone-alt" />
        <code>phone | icon=&quot;phone-alt&quot;</code>
      </div>
      <div>
        <FontAwesomeIcon icon="phone" />
        <code>phone | icon=&quot;phone&quot;</code>
      </div>
      <div>
        <FontAwesomeIcon icon="clock" />
        <code>clock | icon=&quot;clock&quot;</code>
      </div>
      <div>
        <FontAwesomeIcon icon="plus-circle" />
        <code>plus circle | icon=&quot;plus-circle&quot;</code>
      </div>
      <div>
        <FontAwesomeIcon icon="play-circle" />
        <code>play circle | icon=&quot;play-circle&quot;</code>
      </div>
      <div>
        <FontAwesomeIcon icon="plus-square" />
        <code>plus square | icon=&quot;plus-square&quot;</code>
      </div>
      <div>
        <FontAwesomeIcon icon="sort" />
        <code>sort | icon=&quot;sort&quot;</code>
      </div>
      <div>
        <FontAwesomeIcon icon="ban" />
        <code>ban | icon=&quot;ban&quot;</code>
      </div>
      <div>
        <FontAwesomeIcon icon="exclamation-circle" />
        <code>exclamation circle | icon=&quot;exclamation-circle&quot;</code>
      </div>
      <div>
        <FontAwesomeIcon icon="envelope" />
        <code>envelope | icon=&quot;envelope&quot;</code>
      </div>
      <div>
        <FontAwesomeIcon icon="external-link-alt" />
        <code>external link (alt) | icon=&quot;external-link-alt&quot;</code>
      </div>
      <div>
        <FontAwesomeIcon icon="minus-square" />
        <code>minus square | icon=&quot;minus-square&quot;</code>
      </div>
      <div>
        <FontAwesomeIcon icon="sync-alt" />
        <code>sync (alt) | icon=&quot;sync-alt&quot;</code>
      </div>
      <div>
        <FontAwesomeIcon icon="spinner" />
        <code>spinner | icon=&quot;spinner&quot;</code>
      </div>
      <div>
        <FontAwesomeIcon icon="caret-down" />
        <code>caret down | icon=&quot;caret-down&quot;</code>
      </div>
      <div>
        <FontAwesomeIcon icon="exclamation" />
        <code>alert | icon=&quot;exclamation&quot;</code>
      </div>
      <div>
        <ExternalLink />
        <code>
          external link | <strong>uses local svg, not FontAwesome</strong>
        </code>
      </div>
      <div>
        <FontAwesomeIcon className="fa-2x" icon={['far', 'user']} />
        <code>user | icon=&quot;user&quot;</code>
      </div>
      <div>
        <FontAwesomeIcon icon="pencil-alt" />
        <code>pencil-alt | icon=&quot;pencil-alt&quot;</code>
      </div>
      <div>
        <FontAwesomeIcon icon="weight-hanging" />
        <code>weight-hanging | icon=&quot;weight-hanging&quot;</code>
      </div>
      <div>
        <FontAwesomeIcon icon="truck-moving" />
        <code>truck-moving | icon=&quot;truck-moving&quot;</code>
      </div>
      <div>
        <FontAwesomeIcon icon="car" />
        <code>car | icon=&quot;car&quot;</code>
      </div>
      <div>
        <FontAwesomeIcon icon="hand-holding-usd" />
        <code>hand-holding-usd | icon=&quot;hand-holding-usd&quot;</code>
      </div>
      <div>
        <FontAwesomeIcon icon="file-image" />
        <code>file-image | icon=&quot;file-image&quot;</code>
      </div>
      <div>
        <FontAwesomeIcon icon="file-pdf" />
        <code>file-pdf | icon=&quot;file-pdf&quot;</code>
      </div>
    </div>
  </div>
);
