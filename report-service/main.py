import functions_framework
from flask import jsonify
import requests
import os

EXPENSE_API_URL = os.environ.get("EXPENSE_API_URL", "http://localhost:8081")
DEBT_API_URL = os.environ.get("DEBT_API_URL", "http://localhost:8000")

@functions_framework.http
def generate_report(request):
    """
    HTTP Cloud Function para generar reportería.
    No persiste datos, es un agregador Stateless de Expense Service y Debt Calculator.
    """
    # 1. Validación
    group_id = request.args.get('group_id')
    if not group_id:
        return jsonify({"error": "Falta el parámetro group_id"}), 400

    try:
        # 2. Requerir Data a Expense Service
        # Usamos timeout restrictivo asumiendo que los servicios internos deben responder rápido
        exp_resp = requests.get(f"{EXPENSE_API_URL}/expenses?group_id={group_id}", timeout=5)
        expenses = exp_resp.json() if exp_resp.status_code == 200 else []

        # 3. Requerir Deudas Netas
        bal_resp = requests.get(f"{DEBT_API_URL}/balances/{group_id}", timeout=5)
        balances = bal_resp.json() if bal_resp.status_code == 200 else {}
        
        # 4. Requerir Plan de liquidación óptimo
        debt_resp = requests.get(f"{DEBT_API_URL}/debts/{group_id}", timeout=5)
        optimal_debts = debt_resp.json() if debt_resp.status_code == 200 else []

        # 5. Agregación de datos y Resumen
        total_spent = sum([float(e.get("amount", 0)) for e in expenses])
        
        report = {
            "summary": {
                "group_id": group_id,
                "total_expenses_registered": len(expenses),
                "gross_total_spent": total_spent
            },
            "balances": balances,
            "optimal_settlement_plan": optimal_debts,
            "expense_historicals": expenses
        }

        return jsonify(report), 200

    except requests.RequestException as e:
        return jsonify({"error": f"Fallo al comunicarse con microservicios subyacentes: {str(e)}"}), 502
    except Exception as e:
        return jsonify({"error": f"Error inesperado genérico: {str(e)}"}), 500
