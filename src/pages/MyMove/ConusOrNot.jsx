import React, { Component } from 'react';
import { connect } from 'react-redux';
import { func, PropTypes } from 'prop-types';

import SelectableCard from 'components/Customer/SelectableCard';
import { setConusStatus } from 'store/onboarding/actions';
import { selectConusStatus } from 'store/onboarding/selectors';
import { CONUS_STATUS } from 'shared/constants';
import SectionWrapper from 'components/Customer/SectionWrapper';

// eslint-disable-next-line react/prefer-stateless-function
export class ConusOrNot extends Component {
  render() {
    const { setLocation, conusStatus } = this.props;
    const oconusCardText = (
      <>
        Starts or ends in Alaska, Hawaii, or International locations
        <hr className="text-white border-0" />
        <strong>MilMove does not support OCONUS moves yet.</strong> Contact your current transportation office to set up
        your move.
      </>
    );

    return (
      <div className="grid-row">
        <div className="grid-col">
          <h1 className="sm-heading">Where are you moving?</h1>
          <SectionWrapper>
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
          </SectionWrapper>
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
    conusStatus: selectConusStatus(state),
  };
  return props;
};

const mapDispatchToProps = {
  setLocation: setConusStatus,
};

export default connect(mapStateToProps, mapDispatchToProps)(ConusOrNot);
