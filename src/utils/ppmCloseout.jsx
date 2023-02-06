import React from 'react';
import { generatePath, Link } from 'react-router-dom';
import moment from 'moment';

import { formatCents, formatCentsTruncateWhole, formatCustomerDate, formatWeight } from 'utils/formatters';
import { expenseTypeLabels, expenseTypes } from 'constants/ppmExpenseTypes';
import { isExpenseComplete, isWeightTicketComplete, isProGearComplete } from 'utils/shipments';

const getW2Address = (address) => {
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
          hideLabel: true,
        },
        { id: 'startingZIP', label: 'Starting ZIP:', value: ppmShipment.actualPickupPostalCode },
        { id: 'endingZIP', label: 'Ending ZIP:', value: ppmShipment.actualDestinationPostalCode },
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
      renderEditLink: () => (editPath ? <Link to={generatePath(editPath, editParams)}>Edit</Link> : ''),
    },
  ];
};

export const formatWeightTicketItems = (weightTickets, editPath, editParams, handleDelete) => {
  return weightTickets?.map((weightTicket, i) => ({
    id: weightTicket.id,
    isComplete: isWeightTicketComplete(weightTicket),
    draftMessage: 'This trip is missing required information.',
    subheading: <h4 className="text-bold">Trip {i + 1}</h4>,
    rows: [
      {
        id: `vehicleDescription-${i}`,
        label: 'Vehicle description:',
        value: weightTicket.vehicleDescription,
        hideLabel: true,
      },
      { id: `emptyWeight-${i}`, label: 'Empty:', value: formatWeight(weightTicket.emptyWeight) },
      { id: `fullWeight-${i}`, label: 'Full:', value: formatWeight(weightTicket.fullWeight) },
      {
        id: `tripWeight-${i}`,
        label: 'Trip weight:',
        value: formatWeight(weightTicket.fullWeight - weightTicket.emptyWeight),
      },
    ],
    onDelete: () => handleDelete('weightTicket', weightTicket.id, weightTicket.eTag, `Trip ${i + 1}`),
    renderEditLink: () => (
      <Link to={generatePath(editPath, { ...editParams, weightTicketId: weightTicket.id })}>Edit</Link>
    ),
  }));
};

export const formatProGearItems = (proGears, editPath, editParams, handleDelete) => {
  return proGears?.map((proGear, i) => {
    const weightValues = proGear.hasWeightTickets
      ? { id: 'weight', label: 'Weight:', value: formatWeight(proGear.weight) }
      : { id: 'constructedWeight', label: 'Constructed weight:', value: formatWeight(proGear.weight) };
    return {
      id: proGear.id,
      isComplete: isProGearComplete(proGear),
      draftMessage: 'This set is missing required information.',
      subheading: <h4 className="text-bold">Set {i + 1}</h4>,
      rows: [
        {
          id: 'proGearType',
          label: 'Pro-gear Type:',
          value: proGear.belongsToSelf ? 'Pro-gear' : 'Spouse pro-gear',
          hideLabel: true,
        },
        { id: 'description', label: 'Description:', value: proGear.description, hideLabel: true },
        weightValues,
      ],
      renderEditLink: () => <Link to={generatePath(editPath, { ...editParams, proGearId: proGear.id })}>Edit</Link>,
      onDelete: () => handleDelete('proGear', proGear.id, proGear.eTag),
    };
  });
};

export const formatExpenseItems = (expenses, editPath, editParams, handleDelete) => {
  return expenses?.map((expense, i) => {
    const contents = {
      id: expense.id,
      isComplete: isExpenseComplete(expense),
      draftMessage: 'This receipt is missing required information.',
      subheading: <h4 className="text-bold">Receipt {i + 1}</h4>,
      rows: [
        {
          id: 'expenseType',
          label: 'Expense Type:',
          value: expenseTypeLabels[expense.movingExpenseType],
          hideLabel: true,
        },
        { id: 'description', label: 'Description:', value: expense.description, hideLabel: true },
        { id: 'amount', label: 'Amount:', value: `$${formatCents(expense.amount)}` },
      ],
      renderEditLink: () => <Link to={generatePath(editPath, { ...editParams, expenseId: expense.id })}>Edit</Link>,
      onDelete: () => handleDelete('expense', expense.id, expense.eTag),
    };

    if (expense.movingExpenseType === expenseTypes.STORAGE) {
      contents.rows.push({
        id: 'daysInStorage',
        label: 'Days in storage:',
        value: 1 + moment(expense.sitEndDate).diff(moment(expense.sitStartDate), 'days'),
      });
    }
    return contents;
  });
};

export const calculateNetWeightForWeightTicket = (weightTicket) => {
  if (
    weightTicket.emptyWeight == null ||
    weightTicket.fullWeight == null ||
    Number.isNaN(Number(weightTicket.emptyWeight)) ||
    Number.isNaN(Number(weightTicket.fullWeight))
  ) {
    return 0;
  }

  return weightTicket.fullWeight - weightTicket.emptyWeight;
};

export const calculateNetWeightForProGearWeightTicket = (weightTicket) => {
  if (weightTicket.weight == null || Number.isNaN(Number(weightTicket.weight))) {
    return 0;
  }

  return weightTicket.weight;
};

export const calculateTotalNetWeightForWeightTickets = (weightTickets = []) => {
  return weightTickets.reduce((prev, curr) => {
    return prev + calculateNetWeightForWeightTicket(curr);
  }, 0);
};

export const calculateTotalNetWeightForProGearWeightTickets = (proGearWeightTickets = []) => {
  return proGearWeightTickets.reduce((prev, curr) => {
    return prev + calculateNetWeightForProGearWeightTicket(curr);
  }, 0);
};

export const calculateTotalMovingExpensesAmount = (movingExpenses = []) => {
  return movingExpenses.reduce((prev, curr) => {
    return curr.amount && !Number.isNaN(Number(curr.amount)) ? prev + curr.amount : prev;
  }, 0);
};
