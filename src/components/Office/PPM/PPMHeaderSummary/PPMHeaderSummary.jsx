import React, { useState } from 'react';
import { string, bool } from 'prop-types';
import classnames from 'classnames';
import { Tag } from '@trussworks/react-uswds';

import HeaderSection, { sectionTypes } from './HeaderSection';
import styles from './PPMHeaderSummary.module.scss';

import LoadingPlaceholder from 'shared/LoadingPlaceholder';
import SomethingWentWrong from 'shared/SomethingWentWrong';
import { usePPMCloseoutQuery } from 'hooks/queries';
import { formatCustomerContactFullAddress } from 'utils/formatters';

const GCCAndIncentiveInfo = ({ ppmShipmentInfo, updatedItemName, setUpdatedItemName, readOnly }) => {
  const { ppmCloseout, isLoading, isError } = usePPMCloseoutQuery(ppmShipmentInfo.id);

  if (isLoading) return <LoadingPlaceholder />;
  if (isError) return <SomethingWentWrong />;

  const incentives = {
    isAdvanceRequested: ppmShipmentInfo.hasRequestedAdvance,
    isAdvanceReceived: ppmShipmentInfo.hasReceivedAdvance,
    advanceAmountRequested: ppmShipmentInfo.advanceAmountRequested,
    advanceAmountReceived: ppmShipmentInfo.advanceAmountReceived,
    grossIncentive: ppmCloseout.grossIncentive + ppmCloseout.SITReimbursement,
    gcc: ppmCloseout.gcc,
    remainingIncentive: ppmCloseout.remainingIncentive + ppmCloseout.SITReimbursement,
  };

  const incentiveFactors = {
    haulType: ppmCloseout.haulType,
    haulPrice: ppmCloseout.haulPrice,
    haulFSC: ppmCloseout.haulFSC,
    packPrice: ppmCloseout.packPrice,
    unpackPrice: ppmCloseout.unpackPrice,
    dop: ppmCloseout.dop,
    ddp: ppmCloseout.ddp,
    sitReimbursement: ppmCloseout.SITReimbursement,
  };

  return (
    <>
      <hr />
      <HeaderSection
        sectionInfo={{
          type: sectionTypes.incentives,
          ...incentives,
        }}
        dataTestId="incentives"
        updatedItemName={updatedItemName}
        setUpdatedItemName={setUpdatedItemName}
        readOnly={readOnly}
      />
      <hr />
      <HeaderSection
        sectionInfo={{ type: sectionTypes.incentiveFactors, ...incentiveFactors }}
        dataTestId="incentiveFactors"
        updatedItemName={updatedItemName}
        setUpdatedItemName={setUpdatedItemName}
        readOnly={readOnly}
      />
    </>
  );
};
export default function PPMHeaderSummary({ ppmShipmentInfo, order, ppmNumber, showAllFields, readOnly }) {
  const [updatedItemName, setUpdatedItemName] = useState('');
  const shipmentInfo = {
    plannedMoveDate: ppmShipmentInfo.expectedDepartureDate,
    actualMoveDate: ppmShipmentInfo.actualMoveDate,
    pickupAddress: ppmShipmentInfo.pickupAddress
      ? formatCustomerContactFullAddress(ppmShipmentInfo.pickupAddress)
      : '—',
    destinationAddress: ppmShipmentInfo.destinationAddress
      ? formatCustomerContactFullAddress(ppmShipmentInfo.destinationAddress)
      : '—',
    pickupAddressObj: ppmShipmentInfo.pickupAddress,
    destinationAddressObj: ppmShipmentInfo.destinationAddress,
    miles: ppmShipmentInfo.miles,
    estimatedWeight: ppmShipmentInfo.estimatedWeight,
    actualWeight: ppmShipmentInfo.actualWeight,
    isActualExpenseReimbursement: ppmShipmentInfo.isActualExpenseReimbursement,
    allowableWeight: ppmShipmentInfo.allowableWeight,
  };

  return (
    <header className={classnames(styles.PPMHeaderSummary)}>
      <div className={styles.header}>
        <h3>PPM {ppmNumber}</h3>
        {shipmentInfo.isActualExpenseReimbursement && (
          <Tag className={styles.tagInfo}>actual expense reimbursement</Tag>
        )}
        <section>
          <HeaderSection
            sectionInfo={{
              type: sectionTypes.shipmentInfo,
              advanceAmountReceived: ppmShipmentInfo.advanceAmountReceived,
              ...shipmentInfo,
            }}
            dataTestId="shipmentInfo"
            updatedItemName={updatedItemName}
            setUpdatedItemName={setUpdatedItemName}
            readOnly={readOnly}
            grade={order?.grade}
            expanded
          />
        </section>
        {showAllFields && (
          <GCCAndIncentiveInfo
            ppmShipmentInfo={ppmShipmentInfo}
            updatedItemName={updatedItemName}
            setUpdatedItemName={setUpdatedItemName}
            readOnly={readOnly}
          />
        )}
      </div>
    </header>
  );
}

PPMHeaderSummary.propTypes = {
  ppmNumber: string.isRequired,
  showAllFields: bool.isRequired,
};

PPMHeaderSummary.defaultProps = {};

// TODO: Add shape/propType/defaults for incentives and GCC components here.
