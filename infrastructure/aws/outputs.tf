output "public_ip" {
  description = "Public IP address"
  value       = aws_eip.app.public_ip
}

output "ssh_command" {
  description = "SSH command to connect"
  value       = "ssh -i ~/.ssh/${var.key_name}.pem ec2-user@${aws_eip.app.public_ip}"
}
