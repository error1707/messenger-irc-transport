package com.transport.irc_messenger

import android.os.Bundle
import androidx.fragment.app.Fragment
import android.view.LayoutInflater
import android.view.View
import android.view.ViewGroup
import androidx.navigation.fragment.findNavController
import com.transport.irc_messenger.databinding.FragmentFirstBinding
import irc_transport.*

/**
 * A simple [Fragment] subclass as the default destination in the navigation.
 */
class FirstFragment : Fragment() {

    private var _binding: FragmentFirstBinding? = null

    // This property is only valid between onCreateView and
    // onDestroyView.
    private val binding get() = _binding!!

    lateinit var chat: IrcTransport
    var contacts = ArrayList<String>()

    override fun onCreateView(
        inflater: LayoutInflater, container: ViewGroup?,
        savedInstanceState: Bundle?
    ): View? {

        _binding = FragmentFirstBinding.inflate(inflater, container, false)
        return binding.root

    }

    override fun onViewCreated(view: View, savedInstanceState: Bundle?) {
        super.onViewCreated(view, savedInstanceState)

        binding.buttonFirst3.setOnClickListener {
            chat = Irc_transport.newIRCTransport(binding.editTextTextPersonName3.text.toString())
        }

        binding.buttonFirst4.setOnClickListener {
            contacts.add(binding.editTextTextPersonName5.text.toString())
            chat.startReceiveMessagesFrom(binding.editTextTextPersonName5.text.toString())
        }

        binding.buttonFirst.setOnClickListener {
            chat.sendMessages(binding.editTextTextPersonName5.text.toString(), binding.editTextTextPersonName4.text.toString())
        }

        binding.buttonFirst2.setOnClickListener {
            try {
                val msg: String = chat.getMessageFrom(binding.editTextTextPersonName5.text.toString())
                binding.editTextTextPersonName.setText(msg)
            } catch (ex: Exception) {
                binding.editTextTextPersonName.setText("No new messages from this user")
            }
        }
    }

    override fun onDestroyView() {
        for (i in contacts) {
            chat.stopReceiveMessagesFrom(i)
        }
        super.onDestroyView()
        _binding = null
    }
}