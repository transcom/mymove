import React, { Component } from 'react';
import { connect } from 'react-redux';
import { func, PropTypes } from 'prop-types';

import SelectableCard from 'components/Customer/SelectableCard';
import { setConusStatus, selectedConusStatus } from 'scenes/Moves/ducks';
import { CONUS_STATUS } from 'shared/constants';

// eslint-disable-next-line react/prefer-stateless-function
export class ConusOrNot extends Component {
  render() {
    const { setLocation, conusStatus } = this.props;
    const oconusCardText = (
      <>
        <p>Starts or ends in Alaska, Hawaii, or International locations</p>
        <p>
          <strong>MilMove does not support OCONUS moves yet.</strong> Contact your current transportation office to set
          up your move.
        </p>
      </>
    );

    return (
      <div className="grid-row">
        <div className="grid-col">
          <h1 className="sm-heading">Where are you moving?</h1>
          <SelectableCard
            id={`input_${CONUS_STATUS.CONUS}`}
            label="CONUS"
            value={CONUS_STATUS.CONUS}
            onChange={(e) => setLocation(e.target.value)}
            name="conusStatus"
            checked={conusStatus === CONUS_STATUS.CONUS}
            cardText="Starts and ends in the continental US"
          />
          <SelectableCard
            id={`input_${CONUS_STATUS.OCONUS}`}
            label="OCONUS"
            value={CONUS_STATUS.OCONUS}
            onChange={(e) => setLocation(e.target.value)}
            name="conusStatus"
            checked={conusStatus === CONUS_STATUS.OCONUS}
            disabled
            cardText={oconusCardText}
          />
        </div>
      </div>
    );
  }
}

ConusOrNot.propTypes = {
  setLocation: func.isRequired,
  conusStatus: PropTypes.string,
};

ConusOrNot.defaultProps = {
  conusStatus: '',
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

export default connect(mapStateToProps, mapDispatchToProps)(ConusOrNot);
