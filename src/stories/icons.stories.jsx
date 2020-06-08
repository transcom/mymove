import React from 'react';
import { storiesOf } from '@storybook/react';

import { ReactComponent as AddIcon } from 'shared/icon/add.svg';
import { ReactComponent as AlertIcon } from 'shared/icon/alert.svg';
import { ReactComponent as ArrowRightIcon } from 'shared/icon/arrow-right.svg';
import { ReactComponent as ArrowLeftIcon } from 'shared/icon/arrow-left.svg';
import { ReactComponent as CheckmarkIcon } from 'shared/icon/checkmark.svg';
import { ReactComponent as ChevLeftIcon } from 'shared/icon/chevron-left.svg';
import { ReactComponent as ChevRightIcon } from 'shared/icon/chevron-right.svg';
import { ReactComponent as DocsIcon } from 'shared/icon/documents.svg';
import { ReactComponent as DocMenu } from 'shared/icon/doc-menu.svg';
import { ReactComponent as EditIcon } from 'shared/icon/edit.svg';
import { ReactComponent as ExtenalLinkIcon } from 'shared/icon/external-link.svg';
import { ReactComponent as FormCheckmarkIcon } from 'shared/icon/form-checkmark.svg';
import { ReactComponent as FormDoubleCaratIcon } from 'shared/icon/form-double-carat.svg';
import { ReactComponent as LockIcon } from 'shared/icon/lock.svg';
import { ReactComponent as MapPinIcon } from 'shared/icon/map-pin.svg';
import { ReactComponent as RotateClockwiseIcon } from 'shared/icon/rotate-clockwise.svg';
import { ReactComponent as RotateCounterClockwiseIcon } from 'shared/icon/rotate-counter-clockwise.svg';
import { ReactComponent as XHeavyIcon } from 'shared/icon/x-heavy.svg';
import { ReactComponent as XLightIcon } from 'shared/icon/x-light.svg';
import { ReactComponent as ZoomInIcon } from 'shared/icon/zoom-in.svg';
import { ReactComponent as ZoomOutIcon } from 'shared/icon/zoom-out.svg';

// Icons

storiesOf('Global|Icons', module).add('all', () => (
  <div style={{ padding: '20px', background: '#f0f0f0' }}>
    <h3>Icons</h3>
    <div id="icons" style={{ display: 'flex', flexWrap: 'wrap' }}>
      <div>
        <AddIcon />
        <code>add</code>
      </div>
      <div>
        <AlertIcon />
        <code>alert</code>
      </div>
      <div>
        <ArrowRightIcon />
        <code>arrow right</code>
      </div>
      <div>
        <ArrowLeftIcon />
        <code>arrow left</code>
      </div>
      <div>
        <ChevLeftIcon />
        <code>chevron left</code>
      </div>
      <div>
        <ChevRightIcon />
        <code>chevron right</code>
      </div>
      <div>
        <CheckmarkIcon />
        <code>checkmark</code>
      </div>
      <div>
        <DocsIcon />
        <code>documents</code>
      </div>
      <div>
        <DocMenu />
        <code>document menu</code>
      </div>
      <div>
        <EditIcon />
        <code>edit</code>
      </div>
      <div>
        <ExtenalLinkIcon />
        <code>external link</code>
      </div>
      <div>
        <FormCheckmarkIcon />
        <code>form checkmark</code>
      </div>
      <div>
        <FormDoubleCaratIcon />
        <code>form double carat</code>
      </div>
      <div>
        <LockIcon />
        <code>lock</code>
      </div>
      <div>
        <MapPinIcon />
        <code>map pin</code>
      </div>
      <div>
        <RotateClockwiseIcon />
        <code>rotate clockwise</code>
      </div>
      <div>
        <RotateCounterClockwiseIcon />
        <code>rotate counter clockwise</code>
      </div>
      <div>
        <XHeavyIcon />
        <code>x heavy</code>
      </div>
      <div>
        <XLightIcon />
        <code>x light</code>
      </div>
      <div>
        <ZoomInIcon />
        <code>zoom in</code>
      </div>
      <div>
        <ZoomOutIcon />
        <code>zoom out</code>
      </div>
    </div>
  </div>
));
