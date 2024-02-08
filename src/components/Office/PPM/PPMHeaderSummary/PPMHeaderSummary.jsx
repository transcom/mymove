import React from 'react';
import { number, bool } from 'prop-types';
import classnames from 'classnames';

import HeaderSection, { sectionTypes } from './HeaderSection';
import styles from './PPMHeaderSummary.module.scss';

import { PPMCloseoutShape } from 'types/shipment';

export default function PPMHeaderSummary({ ppmCloseout, ppmNumber, showAllFields }) {
  const shipmentInfo = {
    plannedMoveDate: ppmCloseout.plannedMoveDate,
    actualMoveDate: ppmCloseout.actualMoveDate,
    actualPickupPostalCode: ppmCloseout.actualPickupPostalCode,
    actualDestinationPostalCode: ppmCloseout.actualDestinationPostalCode,
    miles: ppmCloseout.miles,
    estimatedWeight: ppmCloseout.estimatedWeight,
    actualWeight: ppmCloseout.actualWeight,
  };
  const incentives = {
    isAdvanceRequested: ppmCloseout.advanceRequested,
    isAdvanceReceived: ppmCloseout.advanceReceived,
    advanceAmountRequested: ppmCloseout.advanceAmountRequested,
    advanceAmountReceived: ppmCloseout.aoa,
    grossIncentive: ppmCloseout.grossIncentive,
    gcc: ppmCloseout.gcc,
    remainingIncentive: ppmCloseout.remainingIncentive,
  };
  const gccFactors = {
    haulPrice: ppmCloseout.haulPrice,
    haulFSC: ppmCloseout.haulFSC,
    fullPackUnpackCharge: ppmCloseout.packPrice + ppmCloseout.unpackPrice,
  };
  return (
    <header className={classnames(styles.PPMHeaderSummary)}>
      <div className={styles.header}>
        <h3>PPM {ppmNumber}</h3>
        <section>
          <HeaderSection
            sectionInfo={{
              type: sectionTypes.shipmentInfo,
              ...shipmentInfo,
            }}
          />
        </section>
        <hr />
        {showAllFields && (
          <>
            <HeaderSection
              sectionInfo={{
                type: sectionTypes.incentives,
                ...incentives,
              }}
            />
            <hr />
            <HeaderSection sectionInfo={{ type: sectionTypes.gcc, ...gccFactors }} />
          </>
        )}
      </div>
    </header>
  );
}

PPMHeaderSummary.propTypes = {
  ppmCloseout: PPMCloseoutShape,
  ppmNumber: number.isRequired,
  showAllFields: bool.isRequired,
};

PPMHeaderSummary.defaultProps = {
  ppmCloseout: undefined,
};

// TODO: Add shape/propType/defaults for incentives and GCC components here.
