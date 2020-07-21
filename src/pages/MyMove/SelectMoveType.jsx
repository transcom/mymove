import React from 'react';
import { connect } from 'react-redux';
import { Radio } from '@trussworks/react-uswds';
import { func } from 'prop-types';

import { SHIPMENT_OPTIONS } from 'shared/constants';
import { setSelectedMoveType } from 'scenes/Moves/ducks';

export const SelectMoveType = (props) => (
  <div className="usa-grid">
    <div className="grid-row">
      <div className="grid-col">
        <h1 className="sm-heading">How do you want to move your belongings?</h1>
        <Radio
          id={SHIPMENT_OPTIONS.PPM}
          label="Arrange it all yourself"
          value={SHIPMENT_OPTIONS.PPM}
          name="moveType"
          onChange={(e) => props.setSelectedMoveType(e.target.value)}
          defaultChecked
        />
        <Radio
          id={SHIPMENT_OPTIONS.HHG}
          label="Have professionals pack and move it all"
          value={SHIPMENT_OPTIONS.HHG}
          onChange={(e) => props.setSelectedMoveType(e.target.value)}
          name="moveType"
        />
      </div>
    </div>
  </div>
);

SelectMoveType.propTypes = {
  setSelectedMoveType: func.isRequired,
};

const mapStateToProps = () => ({});

const mapDispatchToProps = {
  setSelectedMoveType,
};
export default connect(mapStateToProps, mapDispatchToProps)(SelectMoveType);
