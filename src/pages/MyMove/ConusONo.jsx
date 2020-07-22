import React, { Component } from 'react';
import { connect } from 'react-redux';
import { Radio } from '@trussworks/react-uswds';
import { func } from 'prop-types';
import { string } from 'yup';

import { setConusStatus, selectedConusStatus } from 'scenes/Moves/ducks';
import { CONUS_STATUS } from 'shared/constants';

// eslint-disable-next-line react/prefer-stateless-function
export class ConusONo extends Component {
  render() {
    const { setLocation, conusStatus } = this.props;

    return (
      <div className="grid-row">
        <div className="grid-col">
          <h1 className="sm-heading">Where are you moving?</h1>
          <p>Are you moving inside or outside the continental US?</p>
          <Radio
            id={CONUS_STATUS.CONUS}
            label="CONUS (continental US)"
            value={CONUS_STATUS.CONUS}
            name="conusStatus"
            onChange={(e) => setLocation(e.target.value)}
            checked={conusStatus === CONUS_STATUS.CONUS}
          />
          <Radio
            id={CONUS_STATUS.OCONUS}
            label="OCONUS (Alaska, Hawaii, international)"
            value={CONUS_STATUS.OCONUS}
            onChange={(e) => setLocation(e.target.value)}
            name="conusStatus"
            checked={conusStatus === CONUS_STATUS.OCONUS}
          />
        </div>
      </div>
    );
  }
}

ConusONo.propTypes = {
  setLocation: func.isRequired,
  conusStatus: string,
};

ConusONo.defaultProps = {
  conusStatus: CONUS_STATUS.CONUS,
};

const mapStateToProps = (state) => {
  const props = {
    conusStatus: selectedConusStatus(state),
  };
  return props;
};

const mapDispatchToProps = {
  setLocation: setConusStatus,
};

export default connect(mapStateToProps, mapDispatchToProps)(ConusONo);
