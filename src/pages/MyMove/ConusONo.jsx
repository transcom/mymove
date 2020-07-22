import React, { Component } from 'react';
import { connect } from 'react-redux';
import { Radio } from '@trussworks/react-uswds';
import { func } from 'prop-types';
import { string } from 'yup';

import { CONUS_STATUS } from 'shared/constants';
import { setConusStatus, selectedConusStatus } from 'scenes/Moves/ducks';

// eslint-disable-next-line react/prefer-stateless-function
class ConusONo extends Component {
  render() {
    const { setLocation, moveLocation } = this.props;

    return (
      <div className="grid-row">
        <div className="grid-col">
          <h1 className="sm-heading">Where are you moving?</h1>
          <p>Are you moving inside or outside the continental US?</p>
          <Radio
            id={CONUS_STATUS.CONUS}
            label="CONUS (continental US)"
            value={CONUS_STATUS.CONUS}
            name="moveLocation"
            onChange={(e) => setLocation(e.target.value)}
            checked={moveLocation === CONUS_STATUS.CONUS}
          />
          <Radio
            id={CONUS_STATUS.OCONUS}
            label="OCONUS (Alaska, Hawaii, international)"
            value={CONUS_STATUS.OCONUS}
            onChange={(e) => setLocation(e.target.value)}
            name="moveLocation"
            checked={moveLocation === CONUS_STATUS.OCONUS}
          />
        </div>
      </div>
    );
  }
}

ConusONo.propTypes = {
  setLocation: func.isRequired,
  moveLocation: string,
};

ConusONo.defaultProps = {
  moveLocation: CONUS_STATUS.CONUS,
};

const mapStateToProps = (state) => {
  const props = {
    moveLocation: selectedConusStatus(state),
  };
  return props;
};

const mapDispatchToProps = {
  setLocation: setConusStatus,
};

export default connect(mapStateToProps, mapDispatchToProps)(ConusONo);
