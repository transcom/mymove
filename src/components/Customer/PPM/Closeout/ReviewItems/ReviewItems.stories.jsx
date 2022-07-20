import React from 'react';
import { Grid, GridContainer, Link } from '@trussworks/react-uswds';
import { action } from '@storybook/addon-actions';
import { v4 as uuidv4 } from 'uuid';

import ReviewItems from 'components/Customer/PPM/Closeout/ReviewItems/ReviewItems';
import SectionWrapper from 'components/Customer/SectionWrapper';
import { formatCents, formatWeight } from 'utils/formatters';

export default {
  title: 'Customer Components / PPM Closeout / Review Items',
  component: ReviewItems,
  decorators: [
    (Story) => (
      <GridContainer>
        <Grid row>
          <Grid col desktop={{ col: 8, offset: 2 }}>
            <SectionWrapper>
              <Story />
            </SectionWrapper>
          </Grid>
        </Grid>
      </GridContainer>
    ),
  ],
};

const Template = (args) => <ReviewItems {...args} />;

export const AboutYourPPM = Template.bind({});
AboutYourPPM.args = {
  heading: <h2>About Your PPM</h2>,
  contents: [
    {
      id: 'about-your-ppm',
      rows: [
        {
          id: 'departureDate',
          label: 'Departure date:',
          value: '01 Jul 2022',
          hideLabel: true,
        },
        {
          id: 'startingZIP',
          label: 'Starting ZIP:',
          value: '90210',
        },
        {
          id: 'endingZIP',
          label: 'Ending ZIP:',
          value: '30813',
        },
        {
          id: 'advance',
          label: 'Advance:',
          value: 'Yes, $5,987',
        },
      ],
      renderEditLink: () => (
        <Link
          onClick={(e) => {
            e.preventDefault();
            action('edit link clicked');
          }}
          href="#"
        >
          Edit
        </Link>
      ), // use react-router-dom Link when not in Storybook
    },
  ],
};

export const EmptyItems = Template.bind({});
EmptyItems.args = {
  heading: (
    <>
      <h3>Pro-gear</h3>
      <span>(0 lbs)</span>
    </>
  ),
  renderAddButton: () => (
    <Link
      className="usa-button usa-button--secondary"
      onClick={(e) => {
        e.preventDefault();
        action('add pro-gear weight link clicked');
      }}
      href="#"
    >
      Add Pro-gear Weight
    </Link> // use react-router-dom Link when not in Storybook
  ),
  emptyMessage: 'No pro-gear weight documented.',
};

export const WeightTickets = Template.bind({});
WeightTickets.args = {
  heading: (
    <>
      <h3>Weight moved</h3>
      <span>(1,525 lbs)</span>
    </>
  ),
  renderAddButton: () => (
    <Link
      className="usa-button usa-button--secondary"
      onClick={(e) => {
        e.preventDefault();
        action('add weight ticket link clicked');
      }}
      href="#"
    >
      Add More Weight
    </Link> // use react-router-dom Link when not in Storybook
  ),
  contents: [
    {
      id: 'trip-1',
      subheading: <h4 className="text-bold">Trip 1</h4>,
      rows: [
        { id: 'vehicleDescription-1', label: 'Vehicle description:', value: 'DMC Delorean', hideLabel: true },
        { id: 'emptyWeight-1', label: 'Empty:', value: formatWeight(2500) },
        { id: 'fullWeight-1', label: 'Full:', value: formatWeight(3500) },
        {
          id: 'tripWeight-1',
          label: 'Trip Weight:',
          value: formatWeight(1000),
        },
      ],
      renderEditLink: () => (
        <Link
          onClick={(e) => {
            e.preventDefault();
            action('edit link clicked');
          }}
          href="#"
        >
          Edit
        </Link> // use react-router-dom Link when not in Storybook
      ),
      onDelete: () => action('delete button clicked'),
    },
    {
      id: 'trip-2',
      subheading: <h4 className="text-bold">Trip 2</h4>,
      rows: [
        { id: 'vehicleDescription-2', label: 'Vehicle description:', value: 'PT Cruiser', hideLabel: true },
        { id: 'emptyWeight-2', label: 'Empty:', value: formatWeight(2725) },
        { id: 'fullWeight-2', label: 'Full:', value: formatWeight(3250) },
        {
          id: 'tripWeight-2',
          label: 'Trip Weight:',
          value: formatWeight(525),
        },
      ],
      renderEditLink: () => (
        <Link
          onClick={(e) => {
            e.preventDefault();
            action('edit link clicked');
          }}
          href="#"
        >
          Edit
        </Link> // use react-router-dom Link when not in Storybook
      ),
      onDelete: () => action('delete button clicked'),
    },
  ],
};

export const ProGear = Template.bind({});
ProGear.args = {
  heading: (
    <>
      <h3>Pro-gear</h3>
      <span>(2,498 lbs)</span>
    </>
  ),
  renderAddButton: () => (
    <Link
      className="usa-button usa-button--secondary"
      onClick={(e) => {
        e.preventDefault();
        action('add pro-gear weight link clicked');
      }}
      href="#"
    >
      Add Pro-gear Weight
    </Link> // use react-router-dom Link when not in Storybook
  ),
  contents: [
    {
      id: uuidv4(),
      subheading: <h4 className="text-bold">Set 1</h4>,
      rows: [
        { id: 'proGearType', label: 'Pro-gear Type:', value: 'Pro-gear', hideLabel: true },
        { id: 'description', label: 'Description:', value: 'Radio equipment', hideLabel: true },
        { id: 'weight', label: 'Weight:', value: formatWeight(1999) },
      ],
      renderEditLink: () => (
        <Link
          onClick={(e) => {
            e.preventDefault();
            action('edit link clicked');
          }}
          href="#"
        >
          Edit
        </Link> // use react-router-dom Link when not in Storybook
      ),
      onDelete: () => action('delete button clicked'),
    },
    {
      id: uuidv4(),
      subheading: <h4 className="text-bold">Set 2</h4>,
      rows: [
        { id: 'proGearType', label: 'Pro-gear Type:', value: 'Spouse pro-gear', hideLabel: true },
        { id: 'description', label: 'Description:', value: 'Training manuals', hideLabel: true },
        { id: 'constructedWeight:', label: 'Constructed weight:', value: formatWeight(499) },
      ],
      renderEditLink: () => (
        <Link
          onClick={(e) => {
            e.preventDefault();
            action('edit link clicked');
          }}
          href="#"
        >
          Edit
        </Link> // use react-router-dom Link when not in Storybook
      ),
      onDelete: () => action('delete button clicked'),
    },
  ],
};

export const Expenses = Template.bind({});
Expenses.args = {
  heading: (
    <>
      <h3>Expenses</h3>
      <span>($1,005.30)</span>
    </>
  ),
  renderAddButton: () => (
    <Link
      className="usa-button usa-button--secondary"
      onClick={(e) => {
        e.preventDefault();
        action('add expenses link clicked');
      }}
      href="#"
    >
      Add Expenses
    </Link> // use react-router-dom Link when not in Storybook
  ),
  contents: [
    {
      id: uuidv4(),
      subheading: <h4 className="text-bold">Receipt 1</h4>,
      rows: [
        { id: 'expenseType', label: 'Expense Type:', value: 'Packing materials', hideLabel: true },
        { id: 'description', label: 'Description:', value: 'Packing peanuts', hideLabel: true },
        { id: 'amount', label: 'Amount:', value: `$${formatCents(12876)}` },
      ],
      renderEditLink: () => (
        <Link
          onClick={(e) => {
            e.preventDefault();
            action('edit link clicked');
          }}
          href="#"
        >
          Edit
        </Link> // use react-router-dom Link when not in Storybook
      ),
      onDelete: () => action('delete button clicked'),
    },
    {
      id: uuidv4(),
      subheading: <h4 className="text-bold">Receipt 2</h4>,
      rows: [
        { id: 'expenseType', label: 'Expense Type:', value: 'Storage', hideLabel: true },
        { id: 'description', label: 'Description:', value: 'Single unit 100ft²', hideLabel: true },
        { id: 'amount', label: 'Amount:', value: `$${formatCents(87654)}` },
        { id: 'daysInStorage', label: 'Days in storage:', value: '90' },
      ],
      renderEditLink: () => (
        <Link
          onClick={(e) => {
            e.preventDefault();
            action('edit link clicked');
          }}
          href="#"
        >
          Edit
        </Link> // use react-router-dom Link when not in Storybook
      ),
      onDelete: () => action('delete button clicked'),
    },
  ],
};
