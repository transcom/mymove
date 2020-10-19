import React from 'react';
import { connect } from 'react-redux';
import { Radio } from '@trussworks/react-uswds';
import { func, PropTypes } from 'prop-types';

import { setConusStatus, selectedConusStatus } from 'scenes/Moves/ducks';
import { CONUS_STATUS } from 'shared/constants';
import SectionWrapper from 'components/Customer/SectionWrapper';

export const ConusOrNot = ({ setLocation, conusStatus }) => {
  return (
    <>
      <h1>Where are you moving?</h1>
      <SectionWrapper>
        <Radio
          id={CONUS_STATUS.CONUS}
          label="CONUS"
          value={CONUS_STATUS.CONUS}
          name="conusStatus"
          onChange={(e) => setLocation(e.target.value)}
          checked={conusStatus === CONUS_STATUS.CONUS}
        />
        <span className="usa-hint" id="conusStatus">
          Starts and ends in the continental US
        </span>
        <Radio
          id={CONUS_STATUS.OCONUS}
          label="OCONUS"
          value={CONUS_STATUS.OCONUS}
          onChange={(e) => setLocation(e.target.value)}
          name="conusStatus"
          checked={conusStatus === CONUS_STATUS.OCONUS}
        />
        <span className="usa-hint" id="conusStatus">
          Starts or ends in Alaska, Hawaii, or International locations
        </span>
        {conusStatus === CONUS_STATUS.OCONUS && (
          <div>
            MilMove does not support OCONUS moves yet. Contact your current transportation office to set up your move.
          </div>
        )}
      </SectionWrapper>
    </>
  );
};

ConusOrNot.propTypes = {
  setLocation: func.isRequired,
  conusStatus: PropTypes.string,
};

ConusOrNot.defaultProps = {
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

export default connect(mapStateToProps, mapDispatchToProps)(ConusOrNot);
