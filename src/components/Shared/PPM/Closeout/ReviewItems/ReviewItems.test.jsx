import React from 'react';
import { render, screen } from '@testing-library/react';
import { Link } from '@trussworks/react-uswds';

import ReviewItems from 'components/Shared/PPM/Closeout/ReviewItems/ReviewItems';

describe('ReviewItems component', () => {
  it('displays a single review item with required fields', () => {
    render(
      <ReviewItems
        heading={<h2>About Your PPM</h2>}
        contents={[
          {
            id: '1',
            rows: [
              { id: 'moveDate', hideLabel: true, label: 'Actual Move Date:', value: '04 Jul 2022' },
              { id: 'depatureDate', label: 'Departure date:', value: '16 Mar 2020' },
              { id: 'pickupPostalCode', label: 'Starting ZIP:', value: '90210' },
            ],
            renderEditLink: () => <Link to="#">Edit</Link>,
          },
        ]}
      />,
    );

    expect(screen.getByRole('heading', { level: 2 })).toHaveTextContent('About Your PPM');
    expect(screen.queryByRole('term', { name: 'Actual Move Date:' })).not.toBeInTheDocument();
    expect(screen.getAllByRole('term')[0]).toHaveTextContent('Departure date:');
    expect(screen.getAllByRole('term')[1]).toHaveTextContent('Starting ZIP:');
    expect(screen.getByText('Edit')).toBeInstanceOf(HTMLAnchorElement);
  });

  it('displays a single review item with delete action', () => {
    render(
      <ReviewItems
        heading={<h2>About Your PPM</h2>}
        contents={[
          {
            id: '1',
            rows: [
              { id: 'moveDate', hideLabel: true, label: 'Actual Move Date:', value: '04 Jul 2022' },
              { id: 'pickupPostalCode', label: 'Starting ZIP:', value: '90210' },
            ],
            renderEditLink: () => <Link to="#">Edit</Link>,
            onDelete: () => {},
          },
        ]}
      />,
    );

    expect(screen.getByRole('button', { name: 'Delete' })).toBeInTheDocument();
  });

  it('render weight moved required labels', () => {
    render(
      <ReviewItems
        heading={<h3>Weight moved</h3>}
        contents={[
          {
            id: '1',
            subheading: <h4>Trip 1</h4>,
            rows: [
              { id: 'vehicleDescription', label: 'Vehicle description:', value: '2022 Honda CR-V Hybrid' },
              { id: 'trailer', label: 'Trailer:', value: 'No' },
            ],
            renderEditLink: () => <Link to="#">Edit</Link>,
            onDelete: () => {},
          },
        ]}
      />,
    );

    expect(screen.getAllByRole('term')[0]).toHaveTextContent('Vehicle description:');
    expect(screen.getAllByRole('definition')[0]).toHaveTextContent('2022 Honda CR-V Hybrid');
    expect(screen.getAllByRole('term')[1]).toHaveTextContent('Trailer:');
  });

  it('render pro-gear required labels', () => {
    render(
      <ReviewItems
        heading={<h3>Pro-gear</h3>}
        contents={[
          {
            id: '1',
            subheading: <h4>Trip 1</h4>,
            rows: [
              { id: 'proGearDescription', label: 'Description:', value: 'Professional equipment' },
              { id: 'proGearWeight', label: 'Weight:', value: '475 lbs' },
            ],
            renderEditLink: () => <Link to="#">Edit</Link>,
            onDelete: () => {},
          },
        ]}
      />,
    );

    expect(screen.getAllByRole('term')[0]).toHaveTextContent('Description:');
    expect(screen.getAllByRole('term')[1]).toHaveTextContent('Weight:');
  });

  it('render expenses required labels', () => {
    render(
      <ReviewItems
        heading={<h3>Expenses</h3>}
        contents={[
          {
            id: '1',
            subheading: <h4>Receipt 1</h4>,
            rows: [
              { id: 'exprenseType', label: 'Type:', value: 'Packing materials' },
              { id: 'expenseDescription', label: 'Description:', value: 'Packing Peanuts' },
            ],
            renderEditLink: () => <Link to="#">Edit</Link>,
            onDelete: () => {},
          },
        ]}
      />,
    );

    expect(screen.getAllByRole('term')[0]).toHaveTextContent('Type:');
    expect(screen.getAllByRole('term')[1]).toHaveTextContent('Description:');
  });

  it('displays the empty message when there are no contents', () => {
    render(<ReviewItems heading={<h2>Weight</h2>} emptyMessage="There are no items to display" />);

    expect(screen.getByText('There are no items to display')).toBeInTheDocument();
  });

  it('displays the a review items with multiple contents and subheadings', () => {
    render(
      <ReviewItems
        heading={<h3>Pro-gear</h3>}
        renderAddButton={() => <Link to="#">Add Pro-gear Weight</Link>}
        contents={[
          {
            id: '1',
            subheading: <h4>Set 1</h4>,
            rows: [
              { id: 'selfProGear', hideLabel: true, label: 'Pro-gear Type:', value: 'Self pro-gear' },
              { id: 'emptyWeight', label: 'Empty Weight:', value: '833 lbs' },
            ],
            renderEditLink: () => <Link to="#">Edit</Link>,
            onDelete: () => {},
          },
          {
            id: '2',
            subheading: <h4>Set 2</h4>,
            rows: [
              { id: 'selfProGear', hideLabel: true, label: 'Pro-gear Type:', value: 'Spouse pro-gear' },
              { id: 'constructedWeight', label: 'Constructed Weight:', value: '475 lbs' },
            ],
            renderEditLink: () => <Link to="#">Edit</Link>,
            onDelete: () => {},
          },
        ]}
      />,
    );

    expect(screen.getByRole('heading', { level: 3 })).toHaveTextContent('Pro-gear');
    expect(screen.getByText('Add Pro-gear Weight')).toBeInstanceOf(HTMLAnchorElement);

    const terms = screen.getAllByRole('term');
    const definitions = screen.getAllByRole('definition');

    expect(screen.getAllByRole('heading', { level: 4 })[0]).toHaveTextContent('Set 1');
    expect(definitions[0]).toHaveTextContent('Self pro-gear');
    expect(terms[0]).toHaveTextContent('Empty Weight:');
    expect(definitions[1]).toHaveTextContent('833 lbs');

    expect(screen.getAllByRole('button', { name: 'Delete' })[0]).toBeInTheDocument();
    expect(screen.getAllByText('Edit')[0]).toBeInstanceOf(HTMLAnchorElement);

    expect(screen.getAllByRole('heading', { level: 4 })[1]).toHaveTextContent('Set 2');
    expect(definitions[2]).toHaveTextContent('Spouse pro-gear');
    expect(terms[1]).toHaveTextContent('Constructed Weight:');
    expect(definitions[3]).toHaveTextContent('475 lbs');

    expect(screen.getAllByRole('button', { name: 'Delete' })[1]).toBeInTheDocument();
    expect(screen.getAllByText('Edit')[1]).toBeInstanceOf(HTMLAnchorElement);
  });
});
