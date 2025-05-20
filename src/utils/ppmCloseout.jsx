import React from 'react';
import { generatePath } from 'react-router-dom';
import moment from 'moment';

import { formatAddress } from './shipmentDisplay';

import { formatCents, formatCentsTruncateWhole, formatCustomerDate, formatWeight } from 'utils/formatters';
import { expenseTypeLabels, expenseTypes } from 'constants/ppmExpenseTypes';
import { isExpenseComplete, isWeightTicketComplete, isProGearComplete } from 'utils/shipments';
import PPMDocumentsStatus from 'constants/ppms';
import { outlineButtonStyle } from 'shared/standardUI/Buttons/ButtonUsa';
import { LinkButton } from 'shared/standardUI/Buttons';

export const getW2Address = (address) => {
  const addressLine1 = address?.streetAddress2
    ? `${address.streetAddress1} ${address.streetAddress2}`
    : address?.streetAddress1;
  const addressLine2 = `${address?.city}, ${address?.state} ${address?.postalCode}`;
  return (
    <>
      <br />
      {addressLine1}
      <br />
      {addressLine2}
    </>
  );
};

export const formatAboutYourPPMItem = (ppmShipment, editPath, editParams) => {
  return [
    {
      id: 'about-your-ppm',
      isComplete: true,
      rows: [
        {
          id: 'departureDate',
          label: 'Departure date:',
          value: formatCustomerDate(ppmShipment.actualMoveDate),
        },
        { id: 'startingZIP', label: 'Starting Address:', value: formatAddress(ppmShipment.pickupAddress) },
        { id: 'endingZIP', label: 'Ending Address:', value: formatAddress(ppmShipment.destinationAddress) },
        {
          id: 'advance',
          label: 'Advance:',
          value: ppmShipment.hasReceivedAdvance
            ? `Yes, $${formatCentsTruncateWhole(ppmShipment.advanceAmountReceived)}`
            : 'No',
        },
        {
          id: 'w2Address',
          label: 'W-2 address:',
          value: getW2Address(ppmShipment.w2Address),
        },
      ],
      renderEditLink: () =>
        editPath ? (
          <LinkButton className={outlineButtonStyle} to={generatePath(editPath, editParams)}>
            Edit
          </LinkButton>
        ) : (
          ''
        ),
    },
  ];
};

export const formatWeightTicketItems = (weightTickets, editPath, editParams, handleDelete) => {
  return weightTickets?.map((weightTicket, i) => {
    const contents = {
      id: weightTicket.id,
      isComplete: isWeightTicketComplete(weightTicket),
      draftMessage: 'This trip is missing required information.',
      subheading: <h4 className="text-bold">Trip {i + 1}</h4>,
      rows: [
        {
          id: `vehicleDescription-${i}`,
          label: 'Vehicle description:',
          value: weightTicket.vehicleDescription,
        },
        { id: `emptyWeight-${i}`, label: 'Empty:', value: formatWeight(weightTicket.emptyWeight) },
        { id: `fullWeight-${i}`, label: 'Full:', value: formatWeight(weightTicket.fullWeight) },
        {
          id: `tripWeight-${i}`,
          label: 'Trip weight:',
          value: formatWeight(weightTicket.fullWeight - weightTicket.emptyWeight),
        },
        {
          id: `trailer-${i}`,
          label: 'Trailer:',
          value: weightTicket.ownsTrailer ? 'Yes' : 'No',
        },
      ],
      onDelete: () => handleDelete('weightTicket', weightTicket.id, weightTicket.eTag, `Trip ${i + 1}`),
      renderEditLink: () => (
        <LinkButton
          className={outlineButtonStyle}
          data-testid={`weightMoved-${i + 1}`}
          to={generatePath(editPath, { ...editParams, weightTicketId: weightTicket.id })}
        >
          Edit
        </LinkButton>
      ),
    };
    if (weightTicket.vehicleDescription === null) {
      contents.rows.splice(0, 1);
      contents.rows.splice(3, 1);
    }
    return contents;
  });
};

export const formatProGearItems = (proGears, editPath, editParams, handleDelete) => {
  return proGears?.map((proGear, i) => {
    const weightValues =
      proGear.hasWeightTickets !== false
        ? { id: 'weight', label: 'Weight:', value: formatWeight(proGear.weight) }
        : { id: 'constructedWeight', label: 'Constructed weight:', value: formatWeight(proGear.weight) };

    const proGearBelongsToSelf = proGear.belongsToSelf === true || proGear.belongsToSelf === null;
    const description = {
      id: 'description',
      label: 'Description:',
      value: proGear.description ? proGear.description : null,
    };

    const contents = {
      id: proGear.id,
      isComplete: isProGearComplete(proGear),
      draftMessage: 'This set is missing required information.',
      subheading: <h4 className="text-bold">Set {i + 1}</h4>,
      rows: [
        {
          id: 'proGearType',
          label: 'Pro-gear Type:',
          value: proGearBelongsToSelf ? 'Pro-gear' : 'Spouse pro-gear',
          hideLabel: true,
        },
        weightValues,
      ],
      renderEditLink: () => (
        <LinkButton
          className={outlineButtonStyle}
          to={generatePath(editPath, { ...editParams, proGearId: proGear.id })}
        >
          Edit
        </LinkButton>
      ),
      onDelete: () => handleDelete('proGear', proGear.id, proGear.eTag, `Set ${i + 1}`),
    };

    if (proGear.description) {
      contents.rows.splice(1, 0, description);
    }
    return contents;
  });
};

export const formatExpenseItems = (expenses, editPath, editParams, handleDelete) => {
  return expenses?.map((expense, i) => {
    const isSmallPackageExpense = expense.movingExpenseType === expenseTypes.SMALL_PACKAGE;
    const expenseType = {
      id: 'expenseType',
      label: 'Type:',
      value: expenseTypeLabels[expense.movingExpenseType],
    };
    const description =
      expense.movingExpenseType !== expenseTypes.SMALL_PACKAGE
        ? { id: 'description', label: 'Description:', value: expense.description }
        : null;

    const contents = {
      id: expense.id,
      isComplete: isExpenseComplete(expense),
      draftMessage: 'This receipt is missing required information.',
      subheading: (
        <h4 className="text-bold">
          {!isSmallPackageExpense ? 'Receipt ' : 'Package '}
          {i + 1}
        </h4>
      ),
      rows: [{ id: 'amount', label: 'Amount:', value: `$${formatCents(expense.amount)}` }],
      renderEditLink: () => (
        <LinkButton
          className={outlineButtonStyle}
          to={generatePath(editPath, { ...editParams, expenseId: expense.id })}
        >
          Edit
        </LinkButton>
      ),
      onDelete: () => handleDelete('expense', expense.id, expense.eTag, `Receipt ${i + 1}`),
    };

    if (expense.movingExpenseType === expenseTypes.STORAGE) {
      contents.rows.push({
        id: 'daysInStorage',
        label: 'Days in storage:',
        value: 1 + moment(expense.sitEndDate).diff(moment(expense.sitStartDate), 'days'),
      });
    }

    if (isSmallPackageExpense) {
      const weight = formatWeight(expense.weightShipped);
      contents.rows.push({
        id: 'shippedWeight',
        label: 'Weight:',
        value: weight,
      });
      contents.rows.push({
        id: 'isProGear',
        label: 'Pro-gear:',
        value: expense.isProGear === true ? 'Yes' : 'No',
      });
      if (expense.isProGear === true) {
        contents.rows.push({
          id: 'proGearBelongsToSelf',
          label: 'Spouse Pro-gear:',
          value: expense.proGearBelongsToSelf === true ? 'No' : 'Yes',
        });
      }
    }

    if (expense.movingExpenseType) {
      contents.rows.splice(0, 0, expenseType);
    }

    if (description) {
      contents.rows.splice(1, 0, description);
    }

    return contents;
  });
};

export const calculateTotalMovingExpensesAmount = (movingExpenses = []) => {
  const excludedExpenseStatuses = [PPMDocumentsStatus.EXCLUDED, PPMDocumentsStatus.REJECTED]; //  EXCLUDED and REJECTED expenses aren't included in the total.
  return movingExpenses.reduce((prev, curr) => {
    return curr.amount && !Number.isNaN(Number(curr.amount)) && !excludedExpenseStatuses.includes(curr.status)
      ? prev + curr.amount
      : prev;
  }, 0);
};

export const getNonProGearWeightSPR = (expenses) => {
  return expenses
    .filter((expense) => expense.isProGear !== true)
    .reduce((total, expense) => total + (expense.weightShipped || 0), 0);
};

export const getProGearWeightSPR = (expenses) => {
  return expenses
    .filter((expense) => expense.isProGear === true)
    .reduce((total, expense) => total + (expense.weightShipped || 0), 0);
};

export const getTotalPackageWeightSPR = (expenses) => {
  return expenses.reduce((total, expense) => total + (expense.weightShipped || 0), 0);
};

export const hasProGearSPR = (expenses) => {
  return expenses.some((expense) => expense.isProGear === true) ? 'Yes' : 'No';
};

export const hasSpouseProGearSPR = (expenses) => {
  return expenses.some((expense) => expense.isProGear === true && expense.proGearBelongsToSelf === false)
    ? 'Yes'
    : 'No';
};
