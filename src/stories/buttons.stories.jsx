import React from 'react';
import PropTypes from 'prop-types';
import { action } from '@storybook/addon-actions';
import { Button } from '@trussworks/react-uswds';

import { EditButton } from '../components/form';

// Buttons
const ButtonGroup = ({ className, disabled }) => (
  <div className={className} style={{ padding: '20px', display: 'flex', flexWrap: 'wrap' }}>
    <Button className="margin-left-1" disabled={disabled} onClick={action('Button clicked')}>
      <span>Button</span>
    </Button>
    <EditButton className="margin-left-1" disabled={disabled} onClick={action('Edit Button clicked')} />
    <Button className="margin-left-1" disabled={disabled} secondary onClick={action('clicked')}>
      <span>Button</span>
    </Button>
    <Button disabled={disabled} className="usa-button--small margin-left-1" onClick={action('clicked')}>
      <span>Button</span>
    </Button>
    <EditButton className="margin-left-1" disabled={disabled} small onClick={action('Small Edit Button clicked')} />
    <Button disabled={disabled} secondary className="usa-button--small margin-left-1" onClick={action('clicked')}>
      <span>Button</span>
    </Button>
    <EditButton
      className="margin-left-1"
      disabled={disabled}
      secondary
      small
      onClick={action('Secondary Small Edit Button clicked')}
    />
    <Button disabled={disabled} className="usa-button--unstyled margin-left-1" onClick={action('clicked')}>
      <span>Button</span>
    </Button>
    <EditButton
      className="margin-left-1"
      disabled={disabled}
      unstyled
      onClick={action('Unstyled Edit Button clicked')}
    />
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

export default {
  title: 'Components/Button',
};

export const Default = () => <ButtonGroup />;
export const Active = () => <ButtonGroup className="active" />;
export const Hover = () => <ButtonGroup className="hover" />;
export const Focus = () => <ButtonGroup className="focus" />;
export const Disabled = () => <ButtonGroup disabled />;
