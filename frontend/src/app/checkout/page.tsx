import CheckoutForm from "@/components/CheckoutForm";
import ProtectedCheckout from "@/components/ProtectedCheckout";

export default function CheckoutPage() {
  return (
    <ProtectedCheckout>
      <CheckoutForm />
    </ProtectedCheckout>
  );
}
