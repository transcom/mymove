/* eslint-disable camelcase */
import React from 'react';
import { connect } from 'react-redux';
import { func, PropTypes } from 'prop-types';

import SelectableCard from 'components/Customer/SelectableCard';
import { setConusStatus, selectedConusStatus } from 'scenes/Moves/ducks';
import { CONUS_STATUS } from 'shared/constants';
import SectionWrapper from 'components/Customer/SectionWrapper';
import ConnectedWizardPage from 'shared/WizardPage/index';
import { WizardPageShape } from 'types/customerShapes';
import { no_op } from 'shared/utils';

export const ConusOrNot = ({ setLocation, conusStatus, wizardProps }) => {
  const oconusCardText = (
    <>
      <p>Starts or ends in Alaska, Hawaii, or International locations</p>
      <p>
        <strong>MilMove does not support OCONUS moves yet.</strong> Contact your current transportation office to set up
        your move.
      </p>
    </>
  );

  const canMoveNext = conusStatus === CONUS_STATUS.CONUS;

  return (
    // eslint-disable-next-line react/jsx-props-no-spreading
    <ConnectedWizardPage handleSubmit={no_op} canMoveNext={canMoveNext} {...wizardProps}>
      <h1>Where are you moving?</h1>
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
    </ConnectedWizardPage>
  );
};

ConusOrNot.propTypes = {
  setLocation: func.isRequired,
  conusStatus: PropTypes.string,
  wizardProps: WizardPageShape.isRequired,
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
