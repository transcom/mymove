import React from 'react';
import { render, screen, waitFor } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { Formik } from 'formik';

import ShipmentAccountingCodes from './ShipmentAccountingCodes';

describe('components/Office/ShipmentAccountingCodes', () => {
  describe('rendering', () => {
    it('renders with minimal props', () => {
      render(
        <Formik initialValues={{}}>
          <ShipmentAccountingCodes />
        </Formik>,
      );
      expect(screen.getByRole('heading', { level: 2 })).toHaveTextContent(/Accounting codes/);
      expect(screen.getByRole('heading', { level: 2 })).toHaveTextContent(/Optional/);
      expect(screen.getByText('No TAC code entered.')).toBeInTheDocument();
      expect(screen.getByText('No SAC code entered.')).toBeInTheDocument();

      expect(screen.getByRole('button', { name: 'Add code' })).toBeInTheDocument();
    });

    it('renders one or multiple TACs and SACs', async () => {
      render(
        <Formik initialValues={{}}>
          <ShipmentAccountingCodes TACs={{ HHG: '1234', NTS: '5678' }} SACs={{ HHG: '000012345' }} />
        </Formik>,
      );

      // Multiple codes for a category don't get checked by default
      expect(screen.getByLabelText('1234 (HHG)')).not.toBeChecked();
      expect(screen.getByLabelText('5678 (NTS)')).not.toBeChecked();

      // Single code for a category gets checked by defalt
      await waitFor(() => expect(screen.getByLabelText('000012345 (HHG)')).toBeChecked());

      expect(screen.getByRole('button', { name: 'Add or edit codes' })).toBeInTheDocument();
    });

    it('doesnt render undefined TAC or SAC values', async () => {
      render(
        <Formik initialValues={{}}>
          <ShipmentAccountingCodes TACs={{ HHG: '1234', NTS: undefined }} />
        </Formik>,
      );

      expect(screen.queryByText('(NTS)')).not.toBeInTheDocument();

      // Single code for a category gets checked by defalt
      await waitFor(() => expect(screen.getByLabelText('1234 (HHG)')).toBeChecked());
    });

    it('applies Formik internal values', () => {
      render(
        <Formik initialValues={{ tacType: 'NTS' }}>
          <ShipmentAccountingCodes TACs={{ HHG: '1234', NTS: '5678' }} />
        </Formik>,
      );

      expect(screen.getByLabelText('1234 (HHG)')).not.toBeChecked();
      expect(screen.getByLabelText('5678 (NTS)')).toBeChecked();
    });

    it('does not render the "optional" text when the appropriate param is passed in', () => {
      render(
        <Formik initialValues={{}}>
          <ShipmentAccountingCodes optional={false} />
        </Formik>,
      );

      expect(screen.getByRole('heading', { level: 2 })).not.toHaveTextContent(/Optional/);
    });
  });

  describe('interactions', () => {
    it('clicking the edit codes button fires an event prop', () => {
      const onEditCodesClick = jest.fn();
      render(
        <Formik initialValues={{}}>
          <ShipmentAccountingCodes onEditCodesClick={onEditCodesClick} />
        </Formik>,
      );

      userEvent.click(screen.getByRole('button', { name: 'Add code' }));
      expect(onEditCodesClick).toHaveBeenCalled();
    });

    it('clicking a code sets the value', async () => {
      render(
        <Formik initialValues={{}}>
          <ShipmentAccountingCodes TACs={{ HHG: '1234', NTS: '5678' }} />
        </Formik>,
      );

      expect(screen.getByLabelText('5678 (NTS)')).not.toBeChecked();
      userEvent.click(screen.getByLabelText('5678 (NTS)'));
      await waitFor(() => {
        expect(screen.getByLabelText('5678 (NTS)')).toBeChecked();
      });
    });

    it('clicking "Clear selection" clears the value', async () => {
      render(
        <Formik initialValues={{}}>
          <ShipmentAccountingCodes TACs={{ HHG: '1234', NTS: '5678' }} />
        </Formik>,
      );

      expect(screen.getByLabelText('5678 (NTS)')).not.toBeChecked();
      userEvent.click(screen.getByLabelText('5678 (NTS)'));
      await waitFor(() => {
        expect(screen.getByLabelText('5678 (NTS)')).toBeChecked();
      });

      userEvent.click(screen.getByRole('button', { name: 'Clear selection' }));
      await waitFor(() => {
        expect(screen.getByLabelText('5678 (NTS)')).not.toBeChecked();
      });
    });

    it('clicking "Clear selection" on a single code clears the value', async () => {
      render(
        <Formik initialValues={{}}>
          <ShipmentAccountingCodes TACs={{ NTS: '5678' }} />
        </Formik>,
      );

      expect(screen.getByLabelText('5678 (NTS)')).toBeChecked();
      userEvent.click(screen.getByRole('button', { name: 'Clear selection' }));
      await waitFor(() => {
        expect(screen.getByLabelText('5678 (NTS)')).not.toBeChecked();
      });
    });
  });
});
