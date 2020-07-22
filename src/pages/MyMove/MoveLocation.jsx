import React, { Component } from 'react';
import { connect } from 'react-redux';
import { Radio } from '@trussworks/react-uswds';
import { func } from 'prop-types';
import { string } from 'yup';

import { MOVE_LOCATION } from 'shared/constants';
import { setMoveLocation, selectedMoveLocation } from 'scenes/Moves/ducks';

// eslint-disable-next-line react/prefer-stateless-function
class MoveLocation extends Component {
  render() {
    const { setLocation, moveLocation } = this.props;

    return (
      <div className="grid-row">
        <div className="grid-col">
          <h1 className="sm-heading">Where are you moving?</h1>
          <p>Are you moving inside or outside the continental US?</p>
          <Radio
            id={MOVE_LOCATION.CONUS}
            label="CONUS (continental US)"
            value={MOVE_LOCATION.CONUS}
            name="moveLocation"
            onChange={(e) => setLocation(e.target.value)}
            checked={moveLocation === MOVE_LOCATION.CONUS}
          />
          <Radio
            id={MOVE_LOCATION.OCONUS}
            label="OCONUS (Alaska, Hawaii, international)"
            value={MOVE_LOCATION.OCONUS}
            onChange={(e) => setLocation(e.target.value)}
            name="moveLocation"
            checked={moveLocation === MOVE_LOCATION.OCONUS}
          />
        </div>
      </div>
    );
  }
}

MoveLocation.propTypes = {
  setLocation: func.isRequired,
  moveLocation: string,
};

MoveLocation.defaultProps = {
  moveLocation: MOVE_LOCATION.CONUS,
};

const mapStateToProps = (state) => {
  const props = {
    moveLocation: selectedMoveLocation(state),
  };
  return props;
};

const mapDispatchToProps = {
  setLocation: setMoveLocation,
};

export default connect(mapStateToProps, mapDispatchToProps)(MoveLocation);
