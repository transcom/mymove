import React from 'react';
import { number, bool } from 'prop-types';
import classnames from 'classnames';

import HeaderSection, { sectionTypes } from './HeaderSection';
import styles from './PPMHeaderSummary.module.scss';

import { PPMShipmentShape } from 'types/shipment';

export default function PPMHeaderSummary({ ppmShipment, ppmNumber, showAllFields }) {
  return (
    <header className={classnames(styles.PPMHeaderSummary)}>
      <div className={styles.header}>
        <h3>PPM {ppmNumber}</h3>
        <section>
          <HeaderSection
            sectionInfo={{
              type: sectionTypes.shipmentInfo,
              ...ppmShipment,
            }}
          />
        </section>
        <hr />
        {showAllFields && (
          <>
            <HeaderSection
              sectionInfo={{
                type: sectionTypes.incentives,
                estimatedIncentive: ppmShipment.estimatedIncentive,
                hasRequestedAdvance: ppmShipment.hasRequestedAdvance,
                hasReceivedAdvance: ppmShipment.hasReceivedAdvance,
                advanceAmountReceived: ppmShipment.advanceAmountReceived,
                ...ppmShipment.incentives,
              }}
            />
            <hr />
            <HeaderSection sectionInfo={{ type: sectionTypes.gcc, ...ppmShipment.gcc }} />
          </>
        )}
      </div>
    </header>
  );
}

PPMHeaderSummary.propTypes = {
  ppmShipment: PPMShipmentShape,
  ppmNumber: number.isRequired,
  showAllFields: bool.isRequired,
};

PPMHeaderSummary.defaultProps = {
  ppmShipment: undefined,
};

// TODO: Add shape/propType/defaults for incentives and GCC components here.
