import React from 'react';
import { render, screen, waitFor } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { Formik } from 'formik';

import AccountingCodes from './AccountingCodes';

describe('components/Office/AccountingCodes', () => {
  describe('rendering', () => {
    it('renders with minimal props', () => {
      render(
        <Formik initialValues={{}}>
          <AccountingCodes />
        </Formik>,
      );
      expect(screen.getByRole('heading', { level: 2 })).toHaveTextContent(/Accounting codes/);
      expect(screen.getByRole('heading', { level: 2 })).toHaveTextContent(/Optional/);
      expect(screen.getByText('No TAC code entered.')).toBeInTheDocument();
      expect(screen.getByText('No SAC code entered.')).toBeInTheDocument();

      expect(screen.getByRole('button', { name: 'Add code' })).toBeInTheDocument();
    });

    it('renders one or multiple TACs and SACs', () => {
      render(
        <Formik initialValues={{}}>
          <AccountingCodes TACs={{ hhg: '1234', nts: '5678' }} SACs={{ hhg: '000012345' }} />
        </Formik>,
      );

      // Multiple codes for a category don't get checked by default
      expect(screen.getByLabelText('1234 (HHG)')).not.toBeChecked();
      expect(screen.getByLabelText('5678 (NTS)')).not.toBeChecked();

      // Single code for a category gets checked by defalt
      expect(screen.getByLabelText('000012345 (HHG)')).toBeChecked();

      expect(screen.getByRole('button', { name: 'Add or edit codes' })).toBeInTheDocument();
    });

    it('applies Formik internal values', () => {
      render(
        <Formik initialValues={{ tac: '5678' }}>
          <AccountingCodes TACs={{ hhg: '1234', nts: '5678' }} SACs={{ hhg: '000012345' }} />
        </Formik>,
      );

      expect(screen.getByLabelText('1234 (HHG)')).not.toBeChecked();
      expect(screen.getByLabelText('5678 (NTS)')).toBeChecked();
    });

    it('does not render the "optional" text when the appropriate param is passed in', () => {
      render(
        <Formik initialValues={{}}>
          <AccountingCodes optional={false} />
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
          <AccountingCodes onEditCodesClick={onEditCodesClick} />
        </Formik>,
      );

      userEvent.click(screen.getByRole('button', { name: 'Add code' }));
      expect(onEditCodesClick).toHaveBeenCalled();
    });

    it('clicking a code sets the value', async () => {
      render(
        <Formik initialValues={{}}>
          <AccountingCodes TACs={{ hhg: '1234', nts: '5678' }} />
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
          <AccountingCodes TACs={{ hhg: '1234', nts: '5678' }} />
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
  });
});
