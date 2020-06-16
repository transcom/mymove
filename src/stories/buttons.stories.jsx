import React from 'react';
import PropTypes from 'prop-types';
import { storiesOf } from '@storybook/react';
import { action } from '@storybook/addon-actions';
import { Button } from '@trussworks/react-uswds';

import { DocsButton, EditButton } from '../components/form';

// Buttons

const ButtonGroup = ({ className, disabled }) => (
  <div className={className} style={{ padding: '20px', display: 'flex', flexWrap: 'wrap' }}>
    <Button disabled={disabled} onClick={action('Button clicked')}>
      <span>Button</span>
    </Button>
    <EditButton disabled={disabled} onClick={action('Edit Button clicked')} />
    <Button disabled={disabled} secondary onClick={action('clicked')}>
      <span>Button</span>
    </Button>
    <Button disabled={disabled} className="usa-button--small" onClick={action('clicked')}>
      <span>Button</span>
    </Button>
    <EditButton disabled={disabled} small onClick={action('Small Edit Button clicked')} />
    <Button disabled={disabled} secondary className="usa-button--small" onClick={action('clicked')}>
      <span>Button</span>
    </Button>
    <EditButton disabled={disabled} secondary small onClick={action('Secondary Small Edit Button clicked')} />
    <Button disabled={disabled} className="usa-button--unstyled" onClick={action('clicked')}>
      <span>Button</span>
    </Button>
    <EditButton disabled={disabled} unstyled onClick={action('Unstyled Edit Button clicked')} />
  </div>
);

ButtonGroup.defaultProps = {
  className: '',
  disabled: false,
};

ButtonGroup.propTypes = {
  className: PropTypes.string,
  disabled: PropTypes.bool,
};

storiesOf('Components|Button', module)
  .add('default', () => <ButtonGroup />)
  .add('active', () => <ButtonGroup className="active" />)
  .add('hover', () => <ButtonGroup className="hover" />)
  .add('focus', () => <ButtonGroup className="focus" />)
  .add('disabled', () => <ButtonGroup disabled />);

storiesOf('Components|Icon Buttons', module)
  .add('docs', () => (
    <div style={{ padding: '20px', display: 'flex', flexWrap: 'wrap' }}>
      <DocsButton label="My Documents" onClick={action('Docs button clicked')} />
      <DocsButton small label="My Documents" onClick={action('Docs button clicked')} />
      <DocsButton secondary label="My Documents" onClick={action('Docs button clicked')} />
      <DocsButton secondary small label="My Documents" onClick={action('Docs button clicked')} />
      <DocsButton unstyled label="My Documents" onClick={action('Docs button clicked')} />
      <DocsButton disabled label="My Documents" onClick={action('Docs button clicked')} />
    </div>
  ))
  .add('edit', () => (
    <div style={{ padding: '20px', display: 'flex', flexWrap: 'wrap' }}>
      <EditButton onClick={action('Edit button clicked')} />
      <EditButton small onClick={action('Edit button clicked')} />
      <EditButton secondary onClick={action('Edit button clicked')} />
      <EditButton secondary small onClick={action('Edit button clicked')} />
      <EditButton unstyled onClick={action('Edit button clicked')} />
      <EditButton disabled onClick={action('Edit button clicked')} />
    </div>
  ));
