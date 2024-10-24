// Helper to render with React Query
const renderWithQueryClient = (ui) => {
  const queryClient = new QueryClient();
  return render(<QueryClientProvider client={queryClient}>{ui}</QueryClientProvider>);
};

describe('QueueTable', () => {
  it('renders without crashing', () => {
    renderWithQueryClient(<QueueTable />);
    expect(screen.getByText('Payment requested')).toBeInTheDocument();
  });
});
