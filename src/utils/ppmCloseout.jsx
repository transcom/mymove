import React from 'react';
import { generatePath, Link } from 'react-router-dom';
import moment from 'moment';

import { formatCents, formatCentsTruncateWhole, formatCustomerDate, formatWeight } from 'utils/formatters';
import { expenseTypeLabels, expenseTypes } from 'constants/ppmExpenseTypes';

export const formatAboutYourPPMItem = (ppmShipment, editPath, editParams) => {
  return [
    {
      id: 'about-your-ppm',
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
            ? `Yes, $${formatCentsTruncateWhole(ppmShipment.advanceAmountRequested)}`
            : 'No',
        },
        {
          id: 'empty',
          label: '',
          value: '',
        },
        {
          id: 'w2Address',
          label: 'W-2 heading:',
          value: 'W-2 address',
          hideLabel: true,
        },
        {
          id: 'addressLine1',
          label: 'address line1:',
          value: ppmShipment.w2Address?.streetAddress2
            ? `${ppmShipment.w2Address.streetAddress1} ${ppmShipment.w2Address.streetAddress2}`
            : ppmShipment.w2Address?.streetAddress1,
          hideLabel: true,
        },
        {
          id: 'addressLine2',
          label: 'address line2:',
          value: `${ppmShipment.w2Address?.city}, ${ppmShipment.w2Address?.state} ${ppmShipment.w2Address?.postalCode}`,
          hideLabel: true,
        },
      ],
      renderEditLink: () => <Link to={generatePath(editPath, editParams)}>Edit</Link>,
    },
  ];
};

export const formatWeightTicketItems = (weightTickets, editPath, editParams, handleDelete) => {
  return weightTickets?.map((weightTicket, i) => ({
    id: weightTicket.id,
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
    onDelete: () => handleDelete('weightTicket', weightTicket.id, weightTicket.eTag),
    renderEditLink: () => (
      <Link to={generatePath(editPath, { ...editParams, weightTicketId: weightTicket.id })}>Edit</Link>
    ),
  }));
};

export const formatProGearItems = (proGears, editPath, editParams, handleDelete) => {
  return proGears?.map((proGear, i) => {
    const weightValues = proGear.hasWeightTickets
      ? { id: 'weight', label: 'Weight:', value: formatWeight(proGear.fullWeight - proGear.emptyWeight) }
      : { id: 'constructedWeight', label: 'Constructed weight:', value: formatWeight(proGear.constructedWeight) };
    return {
      id: proGear.id,
      subheading: <h4 className="text-bold">Set {i + 1}</h4>,
      rows: [
        {
          id: 'proGearType',
          label: 'Pro-gear Type:',
          value: proGear.selfProGear ? 'Pro-gear' : 'Spouse pro-gear',
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

export const calculateTotalNetWeightForWeightTickets = (weightTickets = []) => {
  return weightTickets.reduce((prev, curr) => {
    return prev + calculateNetWeightForWeightTicket(curr);
  }, 0);
};

export const calculateNetWeightForProGearWeightTicket = (proGearWeightTicket) => {
  if (proGearWeightTicket.constructedWeight != null && !Number.isNaN(Number(proGearWeightTicket.constructedWeight))) {
    return proGearWeightTicket.constructedWeight;
  }

  if (
    proGearWeightTicket.emptyWeight == null ||
    proGearWeightTicket.fullWeight == null ||
    Number.isNaN(Number(proGearWeightTicket.emptyWeight)) ||
    Number.isNaN(Number(proGearWeightTicket.fullWeight))
  ) {
    return 0;
  }

  return proGearWeightTicket.fullWeight - proGearWeightTicket.emptyWeight;
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
